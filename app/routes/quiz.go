package routes

import (
	handlers "madaurus/dev/material/app/handlers/quizes"
	"madaurus/dev/material/app/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


func QuizRoute(c *gin.Engine, collection *mongo.Collection, moduleCollection *mongo.Collection) {
	quizRoute := c.Group("/quizes")
	quizRoute.POST("/create", middlewares.Authentication(), handlers.CreateQuiz(collection, moduleCollection))
	quizRoute.PUT("/update", middlewares.Authentication(), handlers.UpdateQuiz(collection))
	quizRoute.DELETE("/delete/:id", middlewares.Authentication(), handlers.DeleteQuiz(collection))
	quizRoute.GET("/get/:id", middlewares.Authentication(), handlers.GetQuiz(collection))
	quizRoute.GET("/module/:id", middlewares.Authentication(), handlers.GetQuizesByModuleId(collection))
	quizRoute.GET("/admin",
		// middlewares.Authentication(),
		handlers.GetQuizesByAdmin(collection))
}