package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/logs"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
)

// @Summary Get a Lecture
// @Description Get a Lecture
// @Tags Lecture
// @Produce json
// @Param lectureId path string true "Lecture ID"
// @Success 200 {object} interfaces.APIResponse
// @Router /lecture [GET]
// @Security Bearer
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Param Authorization header string true "Auth Token"
func GetLecture(collection *mongo.Collection, instance *kafka.KafkaInstance) gin.HandlerFunc {
	return func(g *gin.Context) {
		user, errU := utils.GetUserPayload(g)
		if errU != nil {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_TOKEN})
			logs.Error(errU.Error())
			return
		}
		lectureId, errP := g.Params.Get("lectureId")
		if errP != true {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		lectureObjId, errD := primitive.ObjectIDFromHex(lectureId)
		if errD != nil {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		lecture, err := services.GetTeacherLecture(collection, g.Request.Context(), lectureObjId)
		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		go instance.EvaluationProducer(user, "Lecture", int32(5))
		g.JSON(http.StatusOK, gin.H{"lecture": lecture})

	}
}

// @Summary Create a new Lecture
// @Description Create a new Lecture
// @Tags Lecture
// @Accept json
// @Produce json
// @Param lecture body models.Lecture true "Lecture Object"
// @Success 200 {object} interfaces.APIResponse
// @Param sectionId path string true "Section ID"
// @Router /lecture [post]
// @Security Bearer
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Param Authorization header string true "Auth Token"
func CreateLecture(collection *mongo.Collection, client *mongo.Client, permitApi *permit.Client, instance *kafka.KafkaInstance) gin.HandlerFunc {
	return func(g *gin.Context) {
		value, errU := g.Get("user")
		if errU != true {
			g.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		user := value.(*utils.UserDetails)
		sectionId := g.Query("sectionId")
		if sectionId == "" {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		sectionObj, errD := primitive.ObjectIDFromHex(sectionId)
		if errD != nil {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		var lecture models.Lecture
		err := g.ShouldBind(&lecture)
		if err != nil {
			log.Printf("Error binding lecture: %v", err.Error())
			g.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_BODY})
			return
		}
		lecture.TeacherId = user.ID
		err = services.CreateLecture(collection, g.Request.Context(), lecture, sectionObj, permitApi, client)
		if err != nil {
			log.Printf("Error creating lecture: %v", err.Error())
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.LECTURE_NOT_CREATED})
			return
		}
		instance.ResourceCreatingProducer(user, "Lecture", lecture.Name, kafka.PROMO_NOTIFICATION_TYPE)
		g.JSON(http.StatusCreated, gin.H{"message": "Lecture Created Successfully"})
	}
}

// @Summary Update a Lecture
// @Description Update a Lecture
// @Tags Lecture
// @Accept json
// @Produce json
// @Param lecture body models.Lecture true "Lecture Object"
// @Success 200 {object} interfaces.APIResponse
// @Param lectureId path string true "Lecture ID"
// @Router /lecture [put]
// @Security Bearer
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Param Authorization header string true "Auth Token"
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

// @Summary Delete a Lecture
// @Description Delete a Lecture
// @Tags Lecture
// @Accept json
// @Produce json
// @Param lectureId path string true "Lecture ID"
// @Success 200 {object} interfaces.APIResponse
// @Router /lecture [delete]
// @Security Bearer
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Param Authorization header string true "Auth Token"
func DeleteLecture(collection *mongo.Collection) gin.HandlerFunc {
	return func(g *gin.Context) {
		lectureId, errP := g.Params.Get("lectureId")
		if errP != true {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		lectureObjId, _ := primitive.ObjectIDFromHex(lectureId)
		err := services.DeleteLecture(collection, g.Request.Context(), lectureObjId)
		if err != nil {
			g.JSON(400, gin.H{"message": err.Error()})
			return
		}
		g.JSON(http.StatusOK, gin.H{"message": shared.LECTURE_DELETED})
	}
}
