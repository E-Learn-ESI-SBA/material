package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
)

func FileRouter(e *gin.Engine, collection *mongo.Collection) {
	e.Group("files", middlewares.Authentication())
	e.PUT("/:id", handlers.EditFile(collection))
}
