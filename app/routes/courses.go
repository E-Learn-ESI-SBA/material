package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	handlers "madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/middlewares"
	"madaurus/dev/material/app/shared/iam"
)

func CourseRoute(c *gin.Engine, collection *mongo.Collection, permitApi *permit.Client, client *mongo.Client, instance *kafka.KafkaInstance) {
	courseRoute := c.Group("/courses")
	courseRoute.GET("/", middlewares.Authentication(), handlers.GetCoursesByTeacher(collection))
	courseRoute.GET("/admin", middlewares.Authentication(), middlewares.StaticRBAC(iam.ROLEAdminKey), handlers.GetCoursesByAdmin(collection))
	courseRoute.POST("/", middlewares.Authentication(), handlers.CreateCourse(collection, client, permitApi, instance))
	courseRoute.PUT("/", middlewares.Authentication(), handlers.UpdateCourse(collection))
	courseRoute.DELETE("/:id", middlewares.Authentication(), handlers.DeleteCourse(collection))

}
