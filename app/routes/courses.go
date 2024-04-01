package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
)

func CourseRoute(c *gin.Engine, collection *mongo.Collection) {
	courseRoute := c.Group("/courses")
	courseRoute.GET("/teacher", handlers.GetCoursesByTeacher(collection))
	courseRoute.GET("/admin", handlers.GetCoursesByAdmin(collection))
	courseRoute.POST("/create", handlers.CreateCourse(collection))
	courseRoute.PUT("/update", handlers.UpdateCourse(collection))
	courseRoute.DELETE("/delete/:id", handlers.DeleteCourse(collection))

}
