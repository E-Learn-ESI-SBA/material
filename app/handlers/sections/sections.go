package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/utils"
)

func GetSections(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		courseId, errP := g.Params.Get("courseId")
		if errP != true {
			g.JSON(400, gin.H{"error": "CourseId is required"})

		}
		sections, err := services.GetSectionsByCourse(g.Request.Context(), collection, courseId)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"sections": sections})

	}
}

func GetSectionDetails(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		sectionId, errP := g.Params.Get("sectionId")
		if errP != true {
			g.JSON(400, gin.H{"error": "SectionId is required"})
			return
		}
		section, err := services.GetSectionDetailsById(g.Request.Context(), collection, sectionId)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"section": section})
	}
}

func CreateSection(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		var section models.Section
		user := g.MustGet("user").(utils.UserDetails)

		err := g.BindJSON(&section)
		section.TeacherId = user.ID
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.CreateSection(g.Request.Context(), collection, section)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"message": "Section Created Successfully"})
	}
}

func EditSection(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		var section models.Section
		sectionId, errP := g.Params.Get("sectionId")
		if errP != true {
			g.JSON(400, gin.H{"error": "SectionId is required"})
			return
		}
		err := g.BindJSON(&section)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.EditSection(g.Request.Context(), collection, section, sectionId)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"message": "Section Updated Successfully"})
	}
}
func DeleteSection(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		sectionId, errP := g.Params.Get("sectionId")
		if errP != true {
			g.JSON(400, gin.H{"error": "SectionId is required"})
			return
		}
		err := services.DeleteSection(g.Request.Context(), collection, sectionId)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"message": "Section Deleted Successfully"})
	}
}
