package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers/lecture"
	"madaurus/dev/material/app/middlewares"
)

func LectureRoute(c *gin.Engine, collection *mongo.Collection) {
	group := c.Group("/lecture", middlewares.Authentication())
	group.GET("/:lectureId", lecture.GetLecture(collection))
	group.POST("/", lecture.CreateLecture(collection))
	group.PUT("/", lecture.UpdateLecture(collection))
	group.DELETE("/:lectureId", lecture.DeleteLecture(collection))

}
