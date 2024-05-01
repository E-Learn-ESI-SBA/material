package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
	"time"
)

// @Summary Create Course
// @Description Protected Route used to create a course (chapter)
// @Produce json
// @Tags Courses
// @Accept json
// @Param course body models.Course true "Course Object"
// @Param module query string true "Module ID"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /coursesp [POST]
func CreateCourse(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var course models.Course
		err := c.ShouldBind(&course)
		if err != nil {
			log.Printf("Error in binding the course: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_BODY})
			return
		}
		moduleId := c.Query("module")
		if moduleId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_BODY})
			return
		}
		course.ID = primitive.NewObjectID()
		createdAt := time.Now()
		course.CreatedAt = &createdAt
		id, errD := primitive.ObjectIDFromHex(moduleId)
		course.ModuleId = id
		if errD != nil {
			log.Printf("Error in converting module id: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.SERVER_ERROR})
			return
		}
		err = services.CreateCourse(c.Request.Context(), collection, course)
		if err != nil {
			log.Printf("Error in creating course: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.SERVER_ERROR})
			return
		}

		log.Printf("Course Created Successfully: %v", course.ID.Hex())

		c.JSON(http.StatusCreated, gin.H{"message": shared.CREATE_COURSE})
	}
}

// @Summary Update Course
// @Description Protected Route used to update a course (chapter)
// @Produce json
// @Tags Courses
// @Accept json
// @Param course body models.Course true "Course Object"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Router /courses/update [PUT]
func UpdateCourse(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var course models.Course
		user := c.MustGet("user").(*utils.UserDetails)
		err := c.BindJSON(&course)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updatedAt := time.Now()
		course.UpdatedAt = &updatedAt
		err = services.UpdateCourse(c.Request.Context(), collection, course, user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Course Updated Successfully"})
	}

}

// @Summary Delete Course
// @Description Protected Route used to delete a course (chapter)
// @Produce json
// @Tags Courses
// @Accept json
// @Param id path string true "Course ID"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Router /courses/delete/{id} [DELETE]
func DeleteCourse(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseId, errP := c.Params.Get("id")
		if errP != true {
			c.JSON(400, gin.H{"error": errors.New("course ID is Required")})
			return
		}
		id, errD := primitive.ObjectIDFromHex(courseId)
		if errD != nil {
			c.JSON(400, gin.H{"error": shared.REQUIRED_ID})
			return
		}
		err := services.DeleteCourse(c.Request.Context(), collection, id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": shared.DELETE_COURSE})
	}
}

// @Summary Getting Course By Admin
// @Description Protected Route used to get the courses (chapters) by admin id
// @Produce json
// @Tags Courses
// @Accept json
// @Security ApiKeyAuth
// @Success 200 {object} []models.Course
// @Failure 400 {object} interfaces.APiError
// @Router /courses/admin [GET]
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

// GetCoursesByTeacher godoc
// @Summary Getting Course By teacher
// @Description Protected Route used to get the courses (chapters) by teacher id
// @Produce json
// @Tags Courses
// @Accept json
// @Security ApiKeyAuth
// @Success 200 {object} []models.Course
// @Failure 400 {object} interfaces.APiError
// @Router /courses/teacher [GET]
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
		courses, errCourses := services.GetCoursesByInstructor(ctx, collection, user.ID)
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
