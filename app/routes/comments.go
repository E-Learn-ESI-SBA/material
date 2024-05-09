package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
)

func CommentRoute(c *gin.Engine, collection *mongo.Collection, permitApi *permit.Client, client *mongo.Client, userCollection *mongo.Collection) {
	comment := c.Group("/comments", middlewares.Authentication())
	comment.GET("/:courseId", handlers.GetCourseComments(collection, userCollection))
	comment.POST("/", handlers.CreateComment(collection, userCollection))
	comment.PUT("/", handlers.EditComment(collection))
	comment.DELETE("/:commentId", handlers.DeleteComment(collection))
	comment.POST("/:commentId/replay/:replayId", handlers.ReplayToComment(collection))
	comment.PUT("/:commentId/replay/:replayId", handlers.EditReplay(collection))
	comment.DELETE("/:commentId/replay/:replayId", handlers.DeleteReplay(collection))

}
