package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
)

func LectureRoute(c *gin.Engine, collection *mongo.Collection) {
	group := c.Group("/lecture", middlewares.Authentication())
	group.GET("/:lectureId", handlers.GetLecture(collection))
	group.POST("", handlers.CreateLecture(collection))
	group.PUT("", handlers.UpdateLecture(collection))
	group.DELETE("/:lectureId", handlers.DeleteLecture(collection))

}
