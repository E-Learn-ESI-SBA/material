package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/utils"
)

func CreateCourse(course models.Course, collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utils.Validator(course)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		_, err = collection.InsertOne(c.Request.Context(), course)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Course Created Successfully"})
	}
}

func GetCoursesByAdmin(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := collection.Find(context.TODO(), nil)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var courses []models.Course
		err = cursor.All(c.Request.Context(), &courses)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"courses": courses})
	}
}

func GetCoursesByTeacher(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// g
		cursor, err := collection.Find(context.TODO(), bson.D{{"teacher_id", c.Params.Get("teacher_id")}})
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var courses []models.Course
		err = cursor.All(c.Request.Context(), &courses)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"courses": courses})
	}

}
func GetCourse(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := bson.D{{"_id", c.Params.Get("id")}}
		var course []models.Course
		err := collection.FindOne(context.TODO(), filter).Decode(course)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"course": course})
	}
}
