package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/logs"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
)

func GetStreamVideo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		videoId := c.Param("id")
		log.Printf("Here ")
		videoObjId, errD := primitive.ObjectIDFromHex(videoId)
		if errD != nil {

			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		video, err := services.GetVideo(c.Request.Context(), collection, videoObjId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		videoFile, errV := services.GetVideoFile(video.Url)
		if errV != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": errV.Error()})
			return
		}
		c.Writer.Header().Set("Content-Type", "video/mp4")
		c.Stream(func(w io.Writer) bool {
			buffer := make([]byte, 1024)
			for {
				n, err := videoFile.Read(buffer)
				if err != nil {
					if err == io.EOF {
						return false
					}
				}
				w.Write(buffer[:n])
			}
			return true
		})
	}

}

func GetVideo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		videoId := c.Param("id")
		videoObjId, errD := primitive.ObjectIDFromHex(videoId)
		if errD != nil {

			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		video, err := services.GetDetailVideo(c.Request.Context(), collection, videoObjId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": video})
	}

}
func EditVideo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var video models.Video
		videoId := c.Param("id")
		videoObjId, errD := primitive.ObjectIDFromHex(videoId)
		if errD != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		err := c.ShouldBindJSON(&video)
		video.ID = videoObjId
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_BODY})
			return
		}
		err = services.EditVideo(c.Request.Context(), collection, video)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}

	}
}

func OnCompleteVideo(collection *mongo.Collection, kafkaInstance *kafka.KafkaInstance) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, errU := utils.GetUserPayload(c)
		if errU != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_TOKEN})
			log.Printf("Error while getting user %v", errU.Error())
			logs.Error(errU.Error())
			return
		}
		videoId := c.Param("id")
		videoObjId, errD := primitive.ObjectIDFromHex(videoId)
		if errD != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		evaluationPoint, err := services.OnCompleteVideo(c.Request.Context(), collection, videoObjId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		kafkaInstance.EvaluationProducer(user, "Video", evaluationPoint)
	}
}
