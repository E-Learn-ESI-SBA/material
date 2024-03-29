package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"log"
	"madaurus/dev/material/app/models"
)

func main() {
	k := koanf.New("/")
	err := k.Load(file.Provider("config/secrets/env.yaml"), yaml.Parser())
	if err != nil {
		log.Fatal("Env file not found")
	}
	var uri string = k.String("database_uri")

	client := models.DBHandler(uri)
	
	server := gin.Default()
	err = server.Run(":8080")
	if err != nil {
		log.Fatal("Server not started")
	}
	fmt.Println("Server Running on Port 8080")
}
