package main

import (
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
	k := shared.GetSecrets()
	uri := k.String("database_uri")
	client, err := models.DBHandler(uri)
	if err != nil {
		log.Fatal("Database not connected")

	}
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
	routes.ModuleRoute(server, app.ModuleCollection)
	routes.CourseRoute(server, app.ModuleCollection)
	routes.SectionRouter(server, app.ModuleCollection)
	routes.LectureRoute(server, app.ModuleCollection)
	routes.CommentRoute(server, app.CommentsCollection)
	routes.TransactionRoute(server, client, app.ModuleCollection)
	routes.FileRouter(server, app.ModuleCollection)
	log.Println("Server Running on Port 8080")
	err = server.Run(":8080")
	defer sentry.Flush(2 * time.Second)
	if err != nil {
		log.Fatal("Server not started")
	}
}
