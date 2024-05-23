package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	handlers "madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
	"madaurus/dev/material/app/shared/iam"
)

func SectionRouter(engine *gin.Engine, collection *mongo.Collection, permitApi *permit.Client, client *mongo.Client) {
	section := engine.Group("/section")
	section.POST("/", middlewares.Authentication(), handlers.CreateSection(collection, client, permitApi))
	section.PUT("/", middlewares.Authentication(), handlers.EditSection(collection))
	section.DELETE("/:sectionId", middlewares.Authentication(), handlers.DeleteSection(collection))
	section.GET("/admin", middlewares.Authentication(), middlewares.StaticRBAC(iam.ROLEAdminKey), handlers.GetSectionsByAdmin(collection))
}
