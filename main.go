package main

import (
	"github.com/gin-gonic/gin"
	"madaurus/dev/material/app/server"
)

func main() {
	gin.SetMode(gin.DebugMode)
	server.NewGracefulServer().Run().Wait()

}
