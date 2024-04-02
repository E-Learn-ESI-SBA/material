package main

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/routes"
	"madaurus/dev/material/app/shared"
	"net/http"
	"os"
	"time"
)

func main() {
	k := shared.GetSecrets()

	var uri string = k.String("database_uri")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client := models.DBHandler(uri, ctx)
	app := new(models.Application)
	app.CreateApp(client)
	server := gin.Default()
	err := server.Run(":8080")
	if err != nil {
		log.Fatal("Server not started")
	}
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
	err = os.Setenv("JWT_SECRET", k.String("jwt_secret"))
	if err != nil {
		log.Fatal("JWT_SECRET not set")

	}

	if errSentry != nil {
		log.Fatalf("sentry.Init: %s", errSentry)
	}
	server.Use(sentrygin.New(sentrygin.Options{}))
	routes.ModuleRoute(server, app.ModuleCollection)
	routes.CourseRoute(server, app.CourseCollection)
	routes.SectionRouter(server, app.SectionCollection)
	routes.LectureRoute(server, app.LectureCollection)
	routes.CommentRoute(server, app.CommentsCollection)

	// Defer Functions
	fmt.Println("Server Running on Port 8080")
	defer sentry.Flush(2 * time.Second)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
		cancel()
	}()
}
