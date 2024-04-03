package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers/courses"
	"madaurus/dev/material/app/middlewares"
)

func CourseRoute(c *gin.Engine, collection *mongo.Collection) {
	courseRoute := c.Group("/courses")
	courseRoute.GET("/teacher", middlewares.Authentication(), courses.GetCoursesByTeacher(collection))
	courseRoute.GET("/admin", middlewares.Authentication(), courses.GetCoursesByAdmin(collection))
	courseRoute.POST("/create", middlewares.Authentication(), courses.CreateCourse(collection))
	courseRoute.PUT("/update", middlewares.Authentication(), courses.UpdateCourse(collection))
	courseRoute.DELETE("/delete/:id", middlewares.Authentication(), courses.DeleteCourse(collection))

}
