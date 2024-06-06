package routes

import (
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


func QuizRoute(
	c *gin.Engine, 
	collection *mongo.Collection, 
	moduleCollection *mongo.Collection, 
	submissionsCollection *mongo.Collection, 
	) {
	quizRoute := c.Group("/quizes")
	quizRoute.POST("", middlewares.Authentication(), handlers.CreateQuiz(collection, moduleCollection))
	quizRoute.GET("/teacher", middlewares.Authentication(), handlers.GetManyQuizesByTeacherId(collection))
	quizRoute.PUT("/:id", middlewares.Authentication(), handlers.UpdateQuiz(collection))
	quizRoute.DELETE("/:id", middlewares.Authentication(), handlers.DeleteQuiz(collection))
	quizRoute.GET("/:id", middlewares.Authentication(), handlers.GetQuiz(collection))
	quizRoute.GET("/module/:id", middlewares.Authentication(), handlers.GetQuizesByModuleId(collection))
	quizRoute.GET("/:id/questions", middlewares.Authentication(), handlers.GetQuizQuestions(collection))
	quizRoute.POST("/:id/submit", middlewares.Authentication(), handlers.SubmitQuizAnswers(collection, submissionsCollection))
	quizRoute.GET("/:id/teacher", middlewares.Authentication(), handlers.GetQuizResults(collection, submissionsCollection))
	quizRoute.GET("/student", middlewares.Authentication(), handlers.GetQuizesResultsByStudentId(collection, submissionsCollection))
	quizRoute.GET("/:id/student", middlewares.Authentication(), handlers.GetQuizResultByStudentId(collection, submissionsCollection, moduleCollection))
	quizRoute.GET("/admin", middlewares.Authentication(), handlers.GetQuizesByAdmin(collection))
}