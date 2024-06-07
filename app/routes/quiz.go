package routes

import (
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


func QuizRoute(
	c *gin.Engine, 
	collection *mongo.Collection, 
	moduleCollection *mongo.Collection, 
	submissionsCollection *mongo.Collection, 
	instance *kafka.KafkaInstance,
	) {
	quizRoute := c.Group("/quizes")
	quizRoute.POST("", middlewares.Authentication(), handlers.CreateQuiz(collection, moduleCollection, instance))
	quizRoute.GET("/teacher", middlewares.Authentication(), handlers.GetManyQuizesByTeacherId(collection))
	quizRoute.PUT("/:id", middlewares.Authentication(), handlers.UpdateQuiz(collection))
	quizRoute.DELETE("/:id", middlewares.Authentication(), handlers.DeleteQuiz(collection))
	quizRoute.GET("/:id", middlewares.Authentication(), handlers.GetQuiz(collection, submissionsCollection))
	quizRoute.GET("/module/:id", middlewares.Authentication(), handlers.GetQuizesByModuleId(collection))
	quizRoute.GET("/:id/questions", middlewares.Authentication(), handlers.GetQuizQuestions(collection, submissionsCollection))
	quizRoute.POST("/:id/submit", middlewares.Authentication(), handlers.SubmitQuizAnswers(collection, submissionsCollection, instance))
	quizRoute.GET("/:id/teacher", middlewares.Authentication(), handlers.GetQuizResults(collection, submissionsCollection))
	quizRoute.GET("/:id/teacher/result", middlewares.Authentication(), handlers.GetSubmissionDetails(collection, submissionsCollection, moduleCollection))
	quizRoute.GET("/student/result", middlewares.Authentication(), handlers.GetQuizesResultsByStudentId(collection, submissionsCollection))
	quizRoute.GET("/:id/student/result", middlewares.Authentication(), handlers.GetQuizResultByStudentId(collection, submissionsCollection, moduleCollection))
	quizRoute.GET("/admin", middlewares.Authentication(), handlers.GetQuizesByAdmin(collection))
	quizRoute.GET("/student", middlewares.Authentication(), handlers.GetQuizesByStudentId(collection))
}