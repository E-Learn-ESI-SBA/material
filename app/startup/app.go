package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"github.com/lpernett/godotenv"
	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/models"
	"os"
	"time"
)

var ServerSetting = &interfaces.Server{}

var DatabaseSetting = &interfaces.Database{}

var RedisSetting = &interfaces.Redis{}

var AppSetting = &interfaces.App{}

var SentrySetting = &interfaces.Sentry{}
var KafkaSetting = &interfaces.Kafka{}

var cfg *ini.File

// init the configuration instance
func Setup() (*mongo.Client, *interfaces.Application, *permit.Client, *kafka.KafkaInstance, *interfaces.App) {
	var err error
	cfg, err = ini.Load("config/conf.ini")
	if err != nil {
		log.Fatalf("settinf.Setup,fail to parse 'conf.ini': %v", err)
	}
	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("kafka", KafkaSetting)
	mapTo("sentry", SentrySetting)
	mapTo("kafka", KafkaSetting)
	environment := os.Getenv("environment")
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
		DatabaseSetting.Host = os.Getenv("database_uri")
		KafkaSetting.Host = os.Getenv("KAFKA_BROKER")
		ServerSetting.PDP_SERVER = os.Getenv("PDP_SERVER")
		ServerSetting.PDP_TOKEN = os.Getenv("PERMIT_TOKEN")
		SentrySetting.DNS = os.Getenv("sentry_dsn")
		AppSetting.JwtSecret = os.Getenv("JWT_SECRET")
	} else {
		gin.SetMode(gin.DebugMode)
		log.Printf("Server setting: %v", ServerSetting)
		log.Println("Running in debug mode")
		gin.SetMode(gin.DebugMode)
		_ = godotenv.Load()
	}
	SentrySetting.EnableTracing = true
	SentrySetting.DEBUG = true
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	InitSentry(SentrySetting)
	client, app, err := models.Setup(DatabaseSetting)

	if err != nil {
		log.Fatal("Database not connected")
	}
	PermitConfig := config.NewConfigBuilder(ServerSetting.PDP_TOKEN).WithPdpUrl(ServerSetting.PDP_SERVER).Build()
	Permit := permit.NewPermit(PermitConfig)
	kafka := kafka.KafkaInit(*KafkaSetting)
	return client, app, Permit, kafka, AppSetting
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", err)
	}
}
