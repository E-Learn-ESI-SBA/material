package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	handlers "madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
)

func SectionRouter(engine *gin.Engine, collection *mongo.Collection) {
	section := engine.Group("/section")
	//	section.GET("/all/:courseId", middlewares.Authentication(), handlers.GetSections(collection))
	//	section.GET("/details/:sectionId", middlewares.Authentication(), handlers.GetSectionDetails(collection))
	section.POST("/", middlewares.Authentication(), handlers.CreateSection(collection))
	section.PUT("/", middlewares.Authentication(), handlers.EditSection(collection))
	section.DELETE("/:sectionId", middlewares.Authentication(), handlers.DeleteSection(collection))
	//	section.GET("/notes/:sectionId", middlewares.Authentication(), handlers.GetSectionsByStudent(collection))

}
