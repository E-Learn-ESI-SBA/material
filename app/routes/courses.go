package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	handlers "madaurus/dev/material/app/handlers/courses"
	"madaurus/dev/material/app/middlewares"
)

func CourseRoute(c *gin.Engine, collection *mongo.Collection) {
	courseRoute := c.Group("/courses")
	courseRoute.GET("/", middlewares.Authentication(), handlers.GetCoursesByTeacher(collection))
	courseRoute.GET("/admin", middlewares.Authentication(), handlers.GetCoursesByAdmin(collection))
	courseRoute.POST("/", middlewares.Authentication(), handlers.CreateCourse(collection))
	courseRoute.PUT("/", middlewares.Authentication(), handlers.UpdateCourse(collection))
	courseRoute.DELETE("/:id", middlewares.Authentication(), handlers.DeleteCourse(collection))

}
