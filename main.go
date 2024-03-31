package main

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"log"
	"madaurus/dev/material/app/models"
	"net/http"
	"time"
)

func main() {
	k := koanf.New("/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err := k.Load(file.Provider("config/secrets/env.yaml"), yaml.Parser())
	if err != nil {
		log.Fatal("Env file not found")
	}
	var uri string = k.String("database_uri")

	client := models.DBHandler(uri)
	app := new(models.Application)
	app.CreateApp(client)
	server := gin.Default()
	err = server.Run(":8080")
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
	if errSentry != nil {
		log.Fatalf("sentry.Init: %s", errSentry)
	}
	server.Use(sentrygin.New(sentrygin.Options{}))
	// Defer Functions
	fmt.Println("Server Running on Port 8080")
	defer sentry.Flush(2 * time.Second)
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}
