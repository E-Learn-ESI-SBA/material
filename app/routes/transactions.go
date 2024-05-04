package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/middlewares"
	transactions2 "madaurus/dev/material/app/transactions"
)

func TransactionRoute(e *gin.Engine, client *mongo.Client, fileCollection *mongo.Collection) {
	transactions := e.Group("/transactions", middlewares.Authentication())
	transactions.POST("/files", transactions2.CreateFileTransaction(client, fileCollection))
	transactions.DELETE("/files/:id", transactions2.DeleteFileTransaction(client, fileCollection))
}
