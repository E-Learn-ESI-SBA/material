package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
)

func GetLecture(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		lectureId, errP := g.Params.Get("lectureId")
		if errP != true {
			g.JSON(400, gin.H{"error": "LectureId is required"})
			return
		}
		lecture, err := services.GetTeacherLecture(collection, g.Request.Context(), lectureId)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"lecture": lecture})
	}
}

func CreateLecture(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		var lecture models.Lecture
		err := g.BindJSON(&lecture)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.CreateLecture(collection, g.Request.Context(), lecture)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"message": "Lecture Created Successfully"})
	}
}

func UpdateLecture(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		var lecture models.Lecture
		err := g.BindJSON(&lecture)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.UpdateLecture(collection, g.Request.Context(), lecture)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"message": "Lecture Updated Successfully"})
	}
}
func DeleteLecture(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		lectureId, errP := g.Params.Get("lectureId")
		if errP != true {
			g.JSON(400, gin.H{"error": "LectureId is required"})
			return
		}
		err := services.DeleteLecture(collection, g.Request.Context(), lectureId)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"message": "Lecture Deleted Successfully"})
	}
}
