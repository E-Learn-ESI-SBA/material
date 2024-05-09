package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
)

func VideoRouter(e *gin.Engine, permitApi *permit.Client, collection *mongo.Collection) {
	videos := e.Group("/videos")
	videos.GET("/:id", middlewares.Authentication(), handlers.GetVideo(collection))
	videos.PUT("/:id", middlewares.Authentication(), handlers.EditVideo(collection))

}
