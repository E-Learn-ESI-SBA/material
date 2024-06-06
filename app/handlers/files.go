package handlers

import (
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/logs"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func EditFile(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {

		var file models.Files
		fileId := c.Param("id")
		fileObjectId, errD := primitive.ObjectIDFromHex(fileId)
		if errD != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		err := c.ShouldBindJSON(&file)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_BODY})
			return
		}
		file.ID = fileObjectId
		err = services.EditFile(c.Request.Context(), collection, file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": shared.FILE_UPDATED})
	}

}

// @Summary Get a file by id
// @Description Get a file by id
func GetFileById(collection *mongo.Collection, kafkaInstance *kafka.KafkaInstance) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := utils.GetUserPayload(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_TOKEN})
			logs.Error(err.Error())
			return
		}
		fileId := c.Param("id")
		fileObjectId, errD := primitive.ObjectIDFromHex(fileId)
		if errD != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		objectFile, err := services.GetFileObject(c.Request.Context(), collection, fileObjectId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		dir, _ := utils.GetStorageFile("files")
		filePath := path.Join(dir, objectFile.Url)
		c.Header("Content-Disposition", "attachment; filename="+objectFile.Name)
		c.Header("Content-Type", "application/"+objectFile.Type)
		// Send the file to the client
		c.File(filePath)
		var evaluationPoint int32 = 15
		go kafkaInstance.EvaluationProducer(user, "File", evaluationPoint)
	}
}
