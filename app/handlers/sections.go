package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
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
// @Success 200 {object} []models.Section
// @Security Bearer
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Failure 401 {object} interfaces.APIResponse
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

/*
// @Summary Get Section Details
// @Description Protected Route Get Section Details
// @Accept json
// @Tags Section
// @Security Bearer
// @Param sectionId path string true "Section Id"
// @Success 200 {object} models.Section
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Failure 401 {object} interfaces.APIResponse
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
*/
// @Summary Create Section
// @Description Protected Route Create Section
// @Accept json
// @Tags Section
// @Param object body models.Section true "Section Object"
// @Param courseId query string true "Course Id"
// @Security Bearer
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Failure 401 {object} interfaces.APIResponse
func CreateSection(collection *mongo.Collection, client *mongo.Client, permitApi *permit.Client) gin.HandlerFunc {
	return func(g *gin.Context) {
		var section models.Section
		// Block
		courseId := g.Query("courseId")
		if courseId == "" {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_ID})
			return
		}
		// Block
		value, errU := g.Get("user")
		user := value.(*utils.UserDetails)
		if !errU {
			g.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		// Block
		objectCourseId, errD := primitive.ObjectIDFromHex(courseId)
		if errD != nil {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_ID})
			return
		}
		// Block
		err := g.BindJSON(&section)
		section.TeacherId = user.ID
		//		section.CourseID = objectCourseId
		section.ID = primitive.NewObjectID()
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_BODY})
			return
		}
		err = services.CreateSection(g.Request.Context(), collection, section, objectCourseId, permitApi, client)
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
// @Failure 400 {object} interfaces.APIResponse
// @Router /section [PUT]
func EditSection(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Block
		value, errU := g.Get("user")
		user := value.(*utils.UserDetails)
		if !errU {
			g.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		// Block
		var section models.Section
		sectionId, errP := g.Params.Get("sectionId")
		if errP != true {
			g.JSON(400, gin.H{"error": "SectionId is required"})
			return
		}
		sectionObj, errD := primitive.ObjectIDFromHex(sectionId)
		if errD != nil {
			g.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_ID})
			return
		}
		err := g.BindJSON(&section)
		if err != nil {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.EditSection(g.Request.Context(), collection, section, sectionObj, user.ID)
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
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Failure 401 {object} interfaces.APIResponse
// @Failure 403 {object} interfaces.APIResponse
// @Router /section/:id [DELETE]

func DeleteSection(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Block
		value, errU := g.Get("user")
		user := value.(*utils.UserDetails)
		if !errU {
			g.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		// Block
		sectionId, errP := g.Params.Get("sectionId")
		if errP != true {
			g.JSON(http.StatusNotAcceptable, gin.H{"error": shared.REQUIRED_ID})
			return
		}
		// Block
		sectionObjectId, errD := primitive.ObjectIDFromHex(sectionId)
		if errD != nil {
			g.JSON(http.StatusNotAcceptable, gin.H{"error": shared.INVALID_ID})
			return
		}
		err := services.DeleteSection(g.Request.Context(), collection, sectionObjectId, user.ID)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_DELETE_MODULE})
			return
		}
		g.JSON(http.StatusOK, gin.H{"message": shared.DELETE_MODULE})
	}
}

/*
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
			g.JSON(200, gin.H{"data": data})
		}
	}
*/
func GetSectionsByAdmin(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		data, errP := services.GetSectionsByAdmin(g.Request.Context(), collection)
		if errP != nil {
			g.JSON(http.StatusBadRequest, gin.H{"message": errP.Error()})
			return
		}
		g.JSON(http.StatusOK, gin.H{"data": data})
	}
}
