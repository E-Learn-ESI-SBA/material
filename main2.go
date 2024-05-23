package main

import (
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/routes"
	"madaurus/dev/material/app/utils"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/permit"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
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
	//	k := shared.GetSecrets()
	uri := os.Getenv("database_uri")
	if uri == "" {
		log.Fatal("Database URI not set")
	}

	log.Println("Db Url %v", uri)
	client, err := models.DBHandler(uri)
	if err != nil {
		log.Fatal("Database not connected")

	}
	PermitConfig := config.NewConfigBuilder(os.Getenv("PERMIT_TOKEN")).WithPdpUrl(os.Getenv("PDP_SERVER")).Build()
	Permit := permit.NewPermit(PermitConfig)
	var app models.Application
	app = *models.NewApp(client)

	if app.ModuleCollection == nil {
		log.Fatal("Database not connected")
	}
	server := gin.New()

	if errSentry != nil {
		log.Fatalf("sentry.Init: %s", errSentry)
	}

	server.MaxMultipartMemory = 250 * 1024 * 1024
	// Start Middleware
	server.Use(cors.New(configCors))
	// End Middleware
	routes.ModuleRoute(server, app.ModuleCollection, Permit, client)
	routes.CourseRoute(server, app.ModuleCollection, Permit, client)
	routes.SectionRouter(server, app.ModuleCollection, Permit, client)
	routes.LectureRoute(server, app.ModuleCollection, Permit, client)
	routes.CommentRoute(server, app.CommentsCollection, Permit, client, app.UserCollection)
	routes.TransactionRoute(server, client, app.ModuleCollection, Permit)
	routes.FileRouter(server, app.ModuleCollection, Permit, client)
	routes.VideoRouter(server, app.ModuleCollection, Permit, client)
	routes.QuizRoute(server, app.QuizesCollection, app.ModuleCollection, app.SubmissionsCollection)
	server.GET("/auth", func(c *gin.Context) {
		user := utils.LightUser{
			Group:    "2021-8",
			Role:     "admin",
			Username: "admin",
			ID:       "2",
			Email:    "ameri@gmail.com",
			Avatar:   "",
			Year:     "3",
		}
		token, _ := utils.GenerateToken(user, os.Getenv("JWT_SECRET"))
		c.Request.Header.Set("Authorization", "Bearer "+token)
		c.SetCookie("accessToken", token, 3600, "/", "localhost", false, false)
		c.JSON(200, gin.H{"token": token})
	})
	log.Println("Server Running on Port 8080")
	err = server.Run(":8080")
	defer sentry.Flush(2 * time.Second)
	if err != nil {
		log.Fatal("Server not started")
	}
}