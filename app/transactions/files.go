package transactions

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
	"path"
	"time"
)

// @Summary Create a file
// @Description Create a file
// @Tags Files
// @Params sectionId query string true
func CreateFileTransaction(client *mongo.Client, collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// start  Mongo session
		dir, err := GetStorageFile("files")
		sectionId := c.Query("sectionId")

		sectionObjId, errD := primitive.ObjectIDFromHex(sectionId)
		if errD != nil {
			log.Printf("Error converting course id: %v", errD)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
		}
		var fileObject models.Files
		value, errU := c.Get("user")
		if !errU {
			log.Printf("Error getting user: %v", errU)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_CREATE_FILE})
		}
		user := value.(*utils.UserDetails)
		fileObject.TeacherId = user.ID
		fileObject.Name = c.PostForm("name")
		fileObject.Group = c.PostForm("group")
		fileObject.Type = c.PostForm("type")
		fileObject.CreatedAt = time.Now()
		fileObject.UpdatedAt = fileObject.CreatedAt
		file, errF := c.FormFile("file")
		if errF != nil {
			log.Printf("Error getting file: %v", errF)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.UNABLE_CREATE_FILE})
			return
		}

		session, err := client.StartSession()
		if err != nil {
			log.Printf("Error starting session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_CREATE_FILE})
			return
		}
		defer func() {
			session.EndSession(c.Request.Context())
			return
		}()
		transactionOption := options.Transaction().SetReadPreference(readpref.Primary())
		err = session.StartTransaction(transactionOption)
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			err = session.AbortTransaction(c.Request.Context())
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_CREATE_FILE})
			return

		}
		fileObject.ID = primitive.NewObjectID()
		fileObject.Url = fileObject.ID.Hex() + "-" + time.Now().Format("2006-01-02") + "." + fileObject.Type
		fileURI := path.Join(dir, fileObject.Url)
		log.Printf("This File path %v", fileURI)
		file.Filename = fileURI
		err = c.SaveUploadedFile(file, fileURI)
		if err != nil {
			log.Printf("Error saving file: %v", err)
			err = session.AbortTransaction(c.Request.Context())
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_CREATE_FILE})
			return
		}
		err = services.CreateFileObject(c.Request.Context(), collection, sectionObjId, fileObject)
		if err != nil {
			log.Printf("Error inserting file: %v", err)
			err = session.AbortTransaction(c.Request.Context())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}

		err = session.CommitTransaction(c.Request.Context())
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_CREATE_FILE})
			return

		}
		c.JSON(http.StatusCreated, gin.H{"message": shared.CREATE_FILE})

	}
}

// For the file (the actual file, not the file object in the database), we need to implement the following functions:

// @Summary Delete a file
// @Description Delete a file
// @Tags Files
// @Params id path string true
// @Security Bearer
// @Success 200 {string} string "File deleted"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "File not deleted"
// @Router /transactions/files/{id} [delete]
func DeleteFileTransaction(client *mongo.Client, collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileId := c.Param("id")
		ctx, _ := context.WithTimeout(c.Request.Context(), time.Second*10)
		defer func() {
			ctx.Done()
		}()
		fileObjectId, errD := primitive.ObjectIDFromHex(fileId)
		if errD != nil {
			log.Printf("Error converting file id: %v", errD)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		// select only the file object
		rs, err := services.GetFileObject(ctx, collection, fileObjectId)
		if err != nil {
			log.Printf("Error getting file object: %v", err.Error())
			c.JSON(http.StatusNotFound, gin.H{"message": shared.FILE_NOT_DELETED})
			return
		}

		session, errS := client.StartSession()
		if errS != nil {
			log.Printf("Error starting session: %v", errS.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return

		}

		dir, err := GetStorageFile("files")
		if err != nil {
			log.Printf("Error getting storage file: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return
		}
		transactionOption := options.Transaction().SetReadPreference(readpref.Primary())
		err = session.StartTransaction(transactionOption)
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			err = session.AbortTransaction(c.Request.Context())
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return

		}
		errOF := services.DeleteFileObject(c.Request.Context(), collection, fileObjectId)
		if errOF != nil {
			session.AbortTransaction(ctx)
			log.Printf("Error deleting file object: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.FILE_NOT_DELETED})
			return

		}
		errDF := services.DeleteSavedFile(rs.File.Url, dir)
		if errDF != nil {
			log.Printf("Error deleting file object: %v", errDF.Error())
			session.AbortTransaction(ctx)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.FILE_NOT_DELETED})
			return
		}
		err = session.CommitTransaction(ctx)
		if err != nil {
			log.Printf("Error While Commiting the Transaction: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return
		}
		defer func() {
			session.EndSession(ctx)
			return
		}()
		c.JSON(http.StatusOK, gin.H{"message": shared.FILE_DELETED})
		return

	}
}

/*



db.modules.aggregate([
    {
        "$unwind": "$courses"
    },
{
        "$unwind": "$courses.sections"
    },
{
        "$unwind": "$courses.sections.files"
    },
    {
        "$match": {
            "courses.sections.files._id":ObjectId("6636a46fa59f3297bb0f9577")
        }
    },
    {
        "$replaceRoot": {
            "newRoot": {
                "$mergeObjects": [
                    "$$ROOT",
                    {
            "file": {
              _id: "$courses.sections.files._id",
              name: "$courses.sections.files.name",
              url: "$courses.sections.files.url",
              // Add other desired file fields here
            }
          }
                ]
            }
        }
    },
    {
        "$project": {
			"courses": 0,
        }
    },

])



*/
