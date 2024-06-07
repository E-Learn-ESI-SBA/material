package server

import (
	"context"
	"fmt"
	"log"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/logs"
	"madaurus/dev/material/app/middlewares"
	"madaurus/dev/material/app/routes"
	"madaurus/dev/material/app/startup"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
)

type GracefulServer struct {
	httpServer  *http.Server
	stopping    chan bool
	status      bool //server 状态
	App         *interfaces.Application
	MongoClient *mongo.Client
	Permit      *permit.Client
	Kafka       *kafka.KafkaInstance
	AppSetting  *interfaces.App
	sync.Mutex
}

func NewGracefulServer() *GracefulServer {
	return &GracefulServer{
		stopping: make(chan bool),
		status:   false,
	}
}

func (s *GracefulServer) Run() *GracefulServer {
	s.Lock()
	defer s.Unlock()
	if s.status {
		return s
	}

	s.status = true
	s.onBoot()
	go s.listenSign()
	go s.runServer()
	return s
}

func (s *GracefulServer) listenSign() {
	listener := make(chan os.Signal)
	signal.Notify(listener, syscall.SIGTERM, syscall.SIGINT)
	<-listener
	s.stopping <- true
}

func (s *GracefulServer) runServer() {
	defer func() {
		if err := recover(); err != nil {
			s.stopping <- false
		}
	}()
	gin.ForceConsoleColor()
	engine := gin.New()
	s.initMiddleware(engine)
	routes.InitRoutes(s.App, s.Permit, s.MongoClient, engine, s.Kafka)
	Addr := fmt.Sprintf(":%d", startup.ServerSetting.HttpPort)
	s.httpServer = &http.Server{
		Addr:              Addr,
		Handler:           engine,
		ReadTimeout:       6 * time.Second,
		WriteTimeout:      6 * time.Second,
		ReadHeaderTimeout: 500 * time.Millisecond,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	err := s.httpServer.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}

func (s *GracefulServer) Wait() {
	if ok := <-s.stopping; ok {
		logs.Info("The application is exiting...")
		log.Println("The application is exiting...")
	} else {
		logs.Warn("The application has an exception and is exiting...")
		log.Println("The application has an exception and is exiting...")
	}

	ctx, cancelTimer := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancelTimer()
		s.onShutDown()
	}()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		logs.Error("The application shutdown error")
	}
}

func (s *GracefulServer) initMiddleware(engine *gin.Engine) {
	configCors := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowFiles:      true,
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token", "hx-request", "hx-current-url"},
		MaxAge:          12 * time.Hour,
	}

	engine.Use(gin.Logger())
	engine.Use(middlewares.Time())
	engine.Use(cors.New(configCors))
	
}

func (s *GracefulServer) onBoot() {
	rand.NewSource(time.Now().UnixNano())
	s.MongoClient, s.App, s.Permit, s.Kafka, s.AppSetting = startup.Setup()
	logs.Setup(s.AppSetting)
	log.Printf("Server is running on port: %d", s.App.ModuleCollection.Name())
	kafka.ExampleProducer(s.Kafka.Producer)
	//	go kafka.ExampleConsumer(s.Kafka.Consumer)
	go func() {
		kafka.UserMutationHandler(s.Kafka.Consumer, s.App.UserCollection)
	}()

}

func (s *GracefulServer) onShutDown() {
	ctx := context.TODO()
	s.Kafka.Producer.Close()
	s.Kafka.Consumer.Close()
	err := s.MongoClient.Disconnect(ctx)
	if err != nil {
		ctx.Err()
	}
	os.Exit(0)
	defer ctx.Done()
}
