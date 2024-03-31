package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
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
		updatedAt := time.Now()
		course.UpdatedAt = &updatedAt

		filter := bson.D{{"_id", c.Params.Get("id")}}
		update := bson.D{{"$set", course}}
		_, err = collection.UpdateOne(c.Request.Context(), filter, update)
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
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		claims, errClaim := c.Get("user")
		if errClaim != false {
			c.JSON(400, gin.H{"error": errors.New("user Details Claims not Found")})
			cancel()
			return
		}
		user, err := claims.(*utils.UserDetails)
		if err == true {
			c.JSON(400, gin.H{"error": errors.New("user Details are not Compatible")})
			cancel()
			return
		}
		courses, errCourses := services.GetCoursesByTeacher(ctx, collection, user.ID)
		if errCourses != nil {
			c.JSON(400, gin.H{"error": errCourses.Error()})
			cancel()
			return
		}
		c.JSON(200, gin.H{"courses": courses})
		defer func() {
			cancel()
		}()
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

func GetPublicFilteredModules(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var filterModule interfaces.ModuleFilter
		err := c.BindJSON(&filterModule)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			cancel()
			return
		}
		modules, CursorErr := services.GetModulesByFilter(ctx, collection, filterModule, "public", nil)
		if CursorErr != nil {
			c.JSON(400, gin.H{"error": CursorErr.Error()})
			cancel()
			return

		}
		c.JSON(200, gin.H{"modules": modules})
		defer cancel()
	}
}
