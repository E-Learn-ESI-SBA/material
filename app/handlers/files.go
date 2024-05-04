package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"net/http"
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
