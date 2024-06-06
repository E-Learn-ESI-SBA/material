package routes

import (
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/kafka"

	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoutes(App *interfaces.Application, permitApi *permit.Client, MongoClient *mongo.Client, engine *gin.Engine, instance *kafka.KafkaInstance) {
	ModuleRoute(engine, App.ModuleCollection, permitApi, MongoClient)
	CommentRoute(engine, App.CommentsCollection, permitApi, MongoClient, App.UserCollection)
	CourseRoute(engine, App.ModuleCollection, permitApi, MongoClient, instance)
	SectionRouter(engine, App.ModuleCollection, permitApi, MongoClient)
	LectureRoute(engine, App.ModuleCollection, permitApi, MongoClient, instance)
	TransactionRoute(engine, MongoClient, App.ModuleCollection, permitApi, instance)
	FileRouter(engine, App.ModuleCollection, permitApi, MongoClient, instance)
	VideoRouter(engine, App.ModuleCollection, permitApi, MongoClient, instance)
	QuizRoute(engine, App.QuizesCollection, App.ModuleCollection, App.SubmissionsCollection, instance)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to Madaurus Material Services"})
	})

}
