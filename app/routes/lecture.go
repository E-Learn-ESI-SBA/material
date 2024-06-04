package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/middlewares"
)

func LectureRoute(c *gin.Engine, collection *mongo.Collection, permitApi *permit.Client, client *mongo.Client, instance *kafka.KafkaInstance) {
	group := c.Group("/lecture", middlewares.Authentication())
	group.GET("/:lectureId", handlers.GetLecture(collection, instance))
	group.POST("", handlers.CreateLecture(collection, client, permitApi, instance))
	group.PUT("", handlers.UpdateLecture(collection))
	group.DELETE("/:lectureId", handlers.DeleteLecture(collection))

}
