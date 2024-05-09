package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/middlewares"
	transactions2 "madaurus/dev/material/app/transactions"
)

func TransactionRoute(e *gin.Engine, client *mongo.Client, collection *mongo.Collection, permitApi *permit.Client) {
	transactions := e.Group("/transactions", middlewares.Authentication())
	transactions.POST("/files", middlewares.Authentication(), transactions2.CreateFileTransaction(client, collection))
	transactions.DELETE("/files/:id", middlewares.Authentication(), transactions2.DeleteFileTransaction(client, collection))
	transactions.POST("/videos/:sectionId", middlewares.Authentication(), transactions2.CreateVideo(collection, client, permitApi))
	transactions.DELETE("videos/:id", middlewares.Authentication(), transactions2.DeleteVideo(collection, client))

}
