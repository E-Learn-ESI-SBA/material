package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/utils"
	"time"
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

func UpdateCourse(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var course models.Course

		err := c.BindJSON(&course)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		filter := bson.D{{"_id", c.Params.Get("id")}}
		update := bson.D{{"$set", course}}
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Course Updated Successfully"})
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
		claims, errClaim := c.Get("user")
		if errClaim != false {
			c.JSON(400, gin.H{"error": errors.New("User Details Claims not Found")})
			return

		}
		user, err := claims.(*utils.UserDetails)
		if err == true {
			c.JSON(400, gin.H{"error": errors.New("user Details are not Compatible")})
			return
		}
		cursor, errMongo := collection.Find(context.TODO(), bson.D{{"teacher_id", user.ID}})
		if errMongo != nil {
			c.JSON(400, gin.H{"error": errMongo.Error()})
			return
		}
		var courses []models.Course
		errMongo = cursor.All(c.Request.Context(), &courses)
		if errMongo != nil {
			c.JSON(400, gin.H{"error": errMongo.Error()})
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

func GetFilteredCourse(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var filterCourse interfaces.CourseFilter
		err := c.BindJSON(&filterCourse)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		filter := bson.D{{"year", filterCourse.Year}, {"speciality", filterCourse.Speciality}}
		var courses models.Course
		cursor, collectionErr := collection.Find(ctx, filter)
		if collectionErr != nil {
			c.JSON(400, gin.H{"error": collectionErr.Error()})
			return
		}
		errCursor := cursor.All(ctx, &courses)
		if err != nil {
			c.JSON(400, gin.H{"error": errCursor.Error()})
		}
		defer func() {
			err := cursor.Close(ctx)
			if err != nil {

			}
		}()
		defer cancel()
	}
}
