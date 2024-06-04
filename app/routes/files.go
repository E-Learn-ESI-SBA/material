package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/middlewares"
)

func FileRouter(e *gin.Engine, collection *mongo.Collection, permitApi *permit.Client, client *mongo.Client, kafkaInstance *kafka.KafkaInstance) {
	e.Group("files")
	e.PUT("/:id", middlewares.Authentication(), handlers.EditFile(collection))
	e.GET("/:id", handlers.GetFileById(collection, kafkaInstance))

}
