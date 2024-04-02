package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
)

func CommentRoute(c *gin.Engine, collection *mongo.Collection) {
	comment := c.Group("/comment", middlewares.Authentication())
	comment.GET("/:contentId", handlers.GetCourseComments(collection))
	comment.POST("/", handlers.CreateComment(collection))
	comment.PUT("/", handlers.EditComment(collection))
	comment.DELETE("/:commentId", handlers.DeleteComment(collection))
	comment.POST("/:commentId/replay/:replayId", handlers.ReplayToComment(collection))
	comment.PUT("/:commentId/replay/:replayId", handlers.EditReplay(collection))
	comment.DELETE("/:commentId/replay/:replayId", handlers.DeleteReplay(collection))

}
