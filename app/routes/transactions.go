package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/middlewares"
	transactions2 "madaurus/dev/material/app/transactions"
)

func TransactionRoute(e *gin.Engine, client *mongo.Client, collection *mongo.Collection) {
	transactions := e.Group("/transactions", middlewares.Authentication())
	transactions.POST("/files", transactions2.CreateFileTransaction(client, collection))
	transactions.DELETE("/files/:id", transactions2.DeleteFileTransaction(client, collection))
	transactions.POST("/videos", transactions2.CreateVideo(collection))
	transactions.DELETE("videos/:id", transactions2.DeleteVideo(collection))
}
