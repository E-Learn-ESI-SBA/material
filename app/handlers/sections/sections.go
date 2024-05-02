package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
)

// @Summary Create Module
// @Description Protected Route Get Sections
// @Accept json
// @Tags Section
// @Param courseId path string true "Course Id"
// @Success 200 {object} []models.SectionResponse
// @Security Bearer
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Failure 401 {object} interfaces.APiError
// @Router /section/all/{courseId} [POST]
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

// @Summary Get Section Details
// @Description Protected Route Get Section Details
// @Accept json
// @Tags Section
// @Security Bearer
// @Param sectionId path string true "Section Id"
// @Success 200 {object} models.SectionDetailsResponse
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Failure 401 {object} interfaces.APiError
// @Router /section/details/{sectionId} [GET]
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

// @Summary Create Section
// @Description Protected Route Create Section
// @Accept json
// @Tags Section
// @Param object body models.Section true "Section Object"
// @Param courseId query string true "Course Id"
// @Security Bearer
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Failure 401 {object} interfaces.APiError
func CreateSection(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		var section models.Section
		value, errU := g.Get("user")
		courseId := g.Query("courseId")
		user := value.(*utils.UserDetails)
		if courseId == "" {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_ID})
			return
		}
		objectCourseId, errD := primitive.ObjectIDFromHex(courseId)
		if errD != nil {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_ID})
			return
		}
		if !errU {
			g.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		err := g.BindJSON(&section)
		// Update Section Document
		section.TeacherId = user.ID
		section.CourseID = objectCourseId
		section.ID = primitive.NewObjectID()
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_BODY})
			return
		}
		err = services.CreateSection(g.Request.Context(), collection, section)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_CREATE_SECTION})
			return
		}
		g.JSON(http.StatusCreated, gin.H{"message": shared.CREATE_SECTION})
	}
}

// @Summary Edit Section
// @Description Protected Route Edit Section
// @Accept json
// @Tags Section
// @Param object body models.Section true "Section Object"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Router /section [PUT]
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

// @Summary Delete Section
// @Description Protected Route Delete Section
// @Accept json
// @Tags Section
// @Param sectionId path string true "Section Id"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Failure 401 {object} interfaces.APiError
// @Failure 403 {object} interfaces.APiError
// @Router /section/:id [DELETE]

func DeleteSection(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		sectionId, errP := g.Params.Get("sectionId")
		if errP != true {
			g.JSON(http.StatusBadRequest, gin.H{"error": shared.REQUIRED_ID})
			return
		}
		sectionObjectId, errD := primitive.ObjectIDFromHex(sectionId)
		if errD != nil {
			g.JSON(http.StatusBadRequest, gin.H{"error": shared.INVALID_ID})
			return
		}
		err := services.DeleteSection(g.Request.Context(), collection, sectionObjectId)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(200, gin.H{"message": "Section Deleted Successfully"})
	}
}

func GetSectionsByStudent(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		sectionId, err := g.Params.Get("sectionId")
		value, errU := g.Get("user")
		user := value.(*utils.UserDetails)
		if !errU {
			g.JSON(401, gin.H{"error": "Session Not Found"})

		}
		studentId := user.ID

		if !err {
			g.JSON(400, gin.H{"error": "SectionId is required"})
			return
		}

		data, errP := services.GetSectionFromStudent(g.Request.Context(), collection, sectionId, studentId)
		if errP != nil {
			g.JSON(400, gin.H{"error": errP.Error()})
			return
		}
		g.JSON(200, gin.H{"section": data})
	}
}
