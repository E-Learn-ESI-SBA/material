package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/middlewares"
)

func VideoRouter(e *gin.Engine, collection *mongo.Collection, permitApi *permit.Client, client *mongo.Client, kafkaInstance *kafka.KafkaInstance) {
	videos := e.Group("/videos")
	videos.GET("/stream/:id", middlewares.Authentication(), handlers.GetStreamVideo(collection))
	videos.GET("/:id", handlers.GetVideo(collection))
	videos.PUT("/:id", middlewares.Authentication(), handlers.EditVideo(collection))
	videos.PATCH("/complete/:id", middlewares.Authentication(), handlers.OnCompleteVideo(collection, kafkaInstance))

}
