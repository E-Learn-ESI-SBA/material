package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/middlewares"
)

func TransactionRoutes(c *gin.Engine, client *mongo.Client) {
	c.Group("/transactions", middlewares.Authentication())
	c.DELETE("module", trans)

}
