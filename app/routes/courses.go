package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
)

func CourseRoute(c *gin.Engine, collection *mongo.Collection) {
	courseRoute := c.Group("/courses")
	courseRoute.GET("/teacher", middlewares.Authentication(), handlers.GetCoursesByTeacher(collection))
	courseRoute.GET("/admin", middlewares.Authentication(), handlers.GetCoursesByAdmin(collection))
	courseRoute.POST("/create", middlewares.Authentication(), handlers.CreateCourse(collection))
	courseRoute.PUT("/update", middlewares.Authentication(), handlers.UpdateCourse(collection))
	courseRoute.DELETE("/delete/:id", middlewares.Authentication(), handlers.DeleteCourse(collection))

}
