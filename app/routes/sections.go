package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers/sections"
	"madaurus/dev/material/app/middlewares"
)

func SectionRouter(engine *gin.Engine, collection *mongo.Collection) {
	section := engine.Group("/section")
	section.GET("/:courseId", middlewares.Authentication(), sections.GetSections(collection))
	section.GET("/single/:sectionId", middlewares.Authentication(), sections.GetSectionDetails(collection))
	section.POST("/", middlewares.Authentication(), sections.CreateSection(collection))
	section.PUT("/", middlewares.Authentication(), sections.EditSection(collection))
	section.DELETE("/:sectionId", middlewares.Authentication(), sections.DeleteSection(collection))

}
