package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/kafka"
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
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /courses [POST]
func CreateCourse(collection *mongo.Collection, client *mongo.Client, permitApi *permit.Client, instance *kafka.KafkaInstance) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := utils.GetUserPayload(c)
		var course models.Course
		err := c.ShouldBind(&course)
		if err != nil {
			log.Printf("Error in binding the course: %v", err)
			c.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_BODY})
			return
		}
		moduleId := c.Query("module")
		if moduleId == "" {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": shared.REQUIRED_ID})
			return
		}

		id, errD := primitive.ObjectIDFromHex(moduleId)
		//	course.ModuleId = id
		if errD != nil {
			log.Printf("Error in converting module id: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		err = services.CreateCourse(c.Request.Context(), collection, course, id, permitApi, client)
		if err != nil {
			log.Printf("Error in creating course: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		log.Printf("Course Created Successfully: %v", course.ID.Hex())
		instance.ResourceCreatingProducer(user, "Chapter", course.Name, kafka.PROMO_NOTIFICATION_TYPE)
		c.JSON(http.StatusCreated, gin.H{"message": shared.CREATE_COURSE})
	}
}

// @Summary Update Course
// @Description Protected Route used to update a course (chapter)
// @Produce json
// @Tags Courses
// @Accept json
// @Param course body models.Course true "Course Object"
// @Security ApiKeyAuth
// @Param id path string true "Course ID"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APIResponse
// @Param module query string true "Module ID"
// @Router /courses/update [PUT]
func UpdateCourse(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var course models.Course
		value, errG := c.Get("user")
		if errG != true {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		courseId, errId := c.Params.Get("id")
		if errId != true {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_ID})
			return
		}
		id, errD := primitive.ObjectIDFromHex(courseId)
		if errD != nil {
			log.Printf("Error in converting course id: %v", errD)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_ID})
			return
		}
		user := value.(*utils.UserDetails)
		err := c.ShouldBind(&course)
		if err != nil {
			log.Printf("Error in binding the course: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		course.ID = id
		moduleId := c.Query("module")
		if moduleId == "" {
			log.Printf("Error in converting module id: %v", err)

			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Module ID is Required")})
			return
		}
		ModuleId, errD2 := primitive.ObjectIDFromHex(moduleId)
		if errD2 != nil {
			log.Printf("Error in converting module id: %v", errD2)
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Module ID is Required")})
			return
		}
		err = services.UpdateCourse(c.Request.Context(), collection, course, user.ID, ModuleId)
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
// @Param moduleId query string true "Module ID"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APIResponse
// @Router /courses/delete/{id} [DELETE]
func DeleteCourse(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseId, errP := c.Params.Get("id")
		if errP != true {
			c.JSON(400, gin.H{"error": errors.New("course ID is Required")})
			return
		}
		moduleId := c.Query("moduleId")
		if moduleId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.QUERY_NOT_FOUND})
			return
		}

		id, errD := primitive.ObjectIDFromHex(courseId)
		if errD != nil {
			c.JSON(400, gin.H{"error": shared.REQUIRED_ID})
			return
		}
		moduleObjectID, errD := primitive.ObjectIDFromHex(moduleId)
		if errD != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": shared.QUERY_NOT_FOUND})
			return
		}
		err := services.DeleteCourse(c.Request.Context(), collection, id, moduleObjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
// @Failure 400 {object} interfaces.APIResponse
// @Router /courses/admin [GET]
func GetCoursesByAdmin(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		courses, err := services.GetCoursesByAdmin(c.Request.Context(), collection)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"courses": courses})
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
// @Failure 400 {object} interfaces.APIResponse
// @Router /courses/teacher [GET]
func GetCoursesByTeacher(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		claims, errClaim := c.Get("user")
		if errClaim != false {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("user Details Claims not Found")})
			cancel()
			return
		}
		user, err := claims.(*utils.UserDetails)
		if err == true {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("user Details are not Compatible")})
			cancel()
			return
		}
		courses, errCourses := services.GetCoursesByInstructor(ctx, collection, user.ID)
		if errCourses != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": errCourses.Error()})
			cancel()
			return
		}
		c.JSON(http.StatusOK, gin.H{"courses": courses})
		defer func() {
			cancel()
		}()
	}

}
