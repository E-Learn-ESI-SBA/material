package main

import (
	"github.com/gin-contrib/cors"
	"github.com/lpernett/godotenv"
	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/permit"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/routes"
	"madaurus/dev/material/app/shared"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Madaurus Material services
// @version 1.0
// @description This Service is for managing the material of the Madaurus Platform
// @termsOfService http://swagger.io/terms/
// @contact.name Seif Hanachi
// @contact.url http://www.swagger.io/support
// @contact.email s.hannachi@esi-sba.dz
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	mode := os.Getenv("GIN_MODE")
	if mode != "release" {
		gin.SetMode(gin.DebugMode)
		_ = godotenv.Load()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	configCors := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowFiles:      true,
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token", "hx-request", "hx-current-url"},
		MaxAge:          12 * time.Hour,
	}
	k := shared.GetSecrets()
	uri := k.String("database_uri")
	client, err := models.DBHandler(uri)
	if err != nil {
		log.Fatal("Database not connected")

	}
	/*
		Enable Permit
	*/

	PermitConfig := config.NewConfigBuilder(os.Getenv("PERMIT_TOKEN")).WithPdpUrl(os.Getenv("PDP_SERVER")).Build()
	Permit := permit.NewPermit(PermitConfig)
	var app models.Application
	app = *models.NewApp(client)

	if app.ModuleCollection == nil {
		log.Fatal("Database not connected")
	}
	server := gin.New()
	errSentry := sentry.Init(sentry.ClientOptions{
		Dsn:           k.String("sentry_dsn"),
		EnableTracing: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if _, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {

				}
			}
			return event
		},
		Debug: true,
	})
	err = os.Setenv("JWT_SECRET", k.String("JWT_SECRET"))
	if err != nil {
		log.Fatal("JWT_SECRET not set")

	}

	if errSentry != nil {
		log.Fatalf("sentry.Init: %s", errSentry)
	}
	url := "http://localhost:8080/swagger/doc.json"
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL(url)))
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to Madaurus Material Services"})
	})

	server.MaxMultipartMemory = 250 * 1024 * 1024
	// Start Middleware
	server.Use(cors.New(configCors))
	// End Middleware
	routes.ModuleRoute(server, app.ModuleCollection, Permit, client)
	routes.CourseRoute(server, app.ModuleCollection, Permit, client)
	routes.SectionRouter(server, app.ModuleCollection, Permit, client)
	routes.LectureRoute(server, app.ModuleCollection, Permit, client)
	routes.CommentRoute(server, app.CommentsCollection, Permit, client)
	routes.TransactionRoute(server, client, app.ModuleCollection, Permit)
	routes.FileRouter(server, app.ModuleCollection, Permit, client)
	routes.VideoRouter(server, app.ModuleCollection, Permit, client)
	log.Println("Server Running on Port 8080")
	err = server.Run(":8080")
	defer sentry.Flush(2 * time.Second)
	if err != nil {
		log.Fatal("Server not started")
	}
}
