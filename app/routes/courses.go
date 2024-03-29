package routes

import "github.com/gin-gonic/gin"

func GetCourses(server *gin.Engine) {
	server.GET("/courses", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Courses",
		})
	})
}
