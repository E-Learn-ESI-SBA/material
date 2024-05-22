package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/interfaces"
)

func InitRoutes(App *interfaces.Application, permitApi *permit.Client, MongoClient *mongo.Client, engine *gin.Engine) {
	ModuleRoute(engine, App.ModuleCollection, permitApi, MongoClient)
	CommentRoute(engine, App.CommentsCollection, permitApi, MongoClient, App.UserCollection)
	CourseRoute(engine, App.ModuleCollection, permitApi, MongoClient)
	SectionRouter(engine, App.ModuleCollection, permitApi, MongoClient)
	LectureRoute(engine, App.ModuleCollection, permitApi, MongoClient)
	TransactionRoute(engine, MongoClient, App.ModuleCollection, permitApi)
	FileRouter(engine, App.ModuleCollection, permitApi, MongoClient)
	VideoRouter(engine, App.ModuleCollection, permitApi, MongoClient)
	QuizRoute(engine, App.QuizesCollection, App.ModuleCollection, App.SubmissionsCollection)
	url := "http://localhost:8080/swagger/doc.json"
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL(url)))
	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to Madaurus Material Services"})
	})

}
