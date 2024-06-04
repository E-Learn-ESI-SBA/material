package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
	"path"
	"strconv"
	"time"
)

// @Summary Create a video
// @Description Create a video
// @Tags Videos
// @Params sectionId path string true
// @Params group formData string true
// @Params name formData string true
// @Params video formData file true
// @Success 201 {string} string "Video created"
// @Security Bearer
// @Router /videos/{sectionId} [post]

func CreateVideo(collection *mongo.Collection, client *mongo.Client, permitApi *permit.Client, instance *kafka.KafkaInstance) gin.HandlerFunc {
	return func(c *gin.Context) {
		var video models.Video
		user, errU := utils.GetUserPayload(c)
		if errU != nil {
			c.JSON(http.StatusInternalServerError, errU.Error())
			return
		}

		sectionId := c.Param("sectionId")
		sectionIdObj, errD := primitive.ObjectIDFromHex(sectionId)
		video.TeacherId = user.ID
		video.Groups = c.PostFormArray("groups")
		// Convert score string to integer

		score, errC := strconv.Atoi(c.PostForm("score"))
		if errC != nil {
			score = 0
		}
		video.Score = int32(score)
		video.Name = c.PostForm("name")
		videoFile, errF := c.FormFile("video")
		if videoFile.Size > 250*1024*1024 {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"message": shared.FILE_TOO_LARGE})
			return
		}
		typ := videoFile.Header.Get("Content-Type")
		log.Printf("File type: %v", typ)
		if errF != nil {
			log.Printf("Error getting file: %v", errF)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.FILE_NOT_DELETED})
			return
		}
		video.ID = primitive.NewObjectID()
		video.Url = video.ID.Hex() + "-" + time.Now().Format("2006-01-02") + ".mp4"
		if errD != nil {
			log.Printf("Error , Not a valid section id %v", errD.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		transactionOption := options.Transaction().SetReadPreference(readpref.Primary())

		session, errS := client.StartSession()
		if errS != nil {
			log.Printf("Error starting session: %v", errS)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_TO_CREATE_VIDEO})
			return
		}
		errS = session.StartTransaction(transactionOption)
		if errS != nil {
			log.Printf("Error starting session: %v", errS)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_TO_CREATE_VIDEO})
			return
		}
		err := services.CreateVideo(c.Request.Context(), collection, sectionIdObj, video)
		if err != nil {
			session.AbortTransaction(c.Request.Context())
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		dir, errD := utils.GetStorageFile("videos")
		fileURI := path.Join(dir, video.Url)
		relation := "parent"
		parentResourceType := "sections"
		err = utils.CreateResourceInstance(permitApi, "videos", video.ID.Hex(), &sectionId, &parentResourceType, &relation)
		if err != nil {
			session.AbortTransaction(c.Request.Context())
			log.Printf("Error creating permit io instance: %v", err.Error())
			c.JSON(http.StatusFailedDependency, gin.H{"message": shared.UNABLE_TO_CREATE_VIDEO})
			return
		}
		videoFile.Filename = fileURI
		err = c.SaveUploadedFile(videoFile, fileURI)
		if err != nil {
			session.AbortTransaction(c.Request.Context())
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.UNABLE_TO_CREATE_VIDEO})
			return
		}
		err = session.CommitTransaction(c.Request.Context())
		// create permit io instance
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_TO_CREATE_VIDEO})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": shared.VIDEO_CREATED})
		defer func() {
			session.EndSession(c.Request.Context())
			instance.ResourceCreatingProducer(user, "Video", video.Name, kafka.PROMO_NOTIFICATION_TYPE)
			return
		}()

	}
}

// @Summary Delete a video
// @Description Delete a video
// @Tags Videos
// @Params id path string true
// @Success 200 {string} string "Video deleted"
// @Security Bearer
// @Router /videos/{id} [delete]

func DeleteVideo(collection *mongo.Collection, client *mongo.Client) gin.HandlerFunc {
	return func(context *gin.Context) {
		videoId := context.Param("id")
		videoObjId, errD := primitive.ObjectIDFromHex(videoId)
		if errD != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		session, err := client.StartSession()
		transactionOption := options.Transaction().SetReadPreference(readpref.Primary())
		err = session.StartTransaction(transactionOption)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": shared.VIDEO_NOT_DELETED})
			return
		}
		video, errV := services.GetVideo(context.Request.Context(), collection, videoObjId)
		if errV != nil {
			session.AbortTransaction(context.Request.Context())
			context.JSON(http.StatusNotFound, gin.H{"message": shared.VIDEO_NOT_DELETED})
			return
		}

		err = services.DeleteVideo(context.Request.Context(), collection, videoObjId)
		if err != nil {
			session.AbortTransaction(context.Request.Context())
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.VIDEO_NOT_DELETED})
			return
		}
		err = services.DeletePhysicalVideo(video.Url)
		if err != nil {
			session.AbortTransaction(context.Request.Context())
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.VIDEO_NOT_DELETED})
			return
		}
		context.JSON(http.StatusOK, gin.H{"message": shared.VIDEO_DELETED})
		return
	}
}
