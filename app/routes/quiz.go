package routes

import (
	handlers "madaurus/dev/material/app/handlers/quizes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


func QuizRoute(c *gin.Engine, collection *mongo.Collection) {
	quizRoute := c.Group("/quizes")
	quizRoute.POST("/create",
		// middlewares.Authentication(),
		handlers.CreateQuiz(collection))
	quizRoute.PUT("/update",
		// middlewares.Authentication(),
		handlers.UpdateQuiz(collection))
	quizRoute.DELETE("/delete/:id",
		// middlewares.Authentication(),
		handlers.DeleteQuiz(collection))
	quizRoute.GET("/get/:id",
		// middlewares.Authentication(),
		handlers.GetQuiz(collection))
	quizRoute.GET("/admin",
		// middlewares.Authentication(),
		handlers.GetQuizesByAdmin(collection))
}