package transactions

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
	"os"
	"path"
	"time"
)

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
		studentId := c.Query("section_id")
		if studentId == "" {
			log.Printf("Error getting section id: %v", studentId)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
		}
		user := value.(utils.UserDetails)
		createdAt := time.Now()
		fileObject.TeacherId = user.ID
		fileObject.Name = c.PostForm("name")
		fileObject.Group = c.PostForm("group")
		fileObject.Type = c.PostForm("type")
		fileObject.CreatedAt = &createdAt
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
		fileObject.Url = fileObject.ID.String() + "-" + time.DateOnly + "." + fileObject.Type
		fileURI := path.Join(dir, fileObject.Url)
		err = c.SaveUploadedFile(file, fileURI)
		if err != nil {
			log.Printf("Error saving file: %v", err)
			err = session.AbortTransaction(c.Request.Context())
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_CREATE_FILE})
			return
		}
		// update module.courses[wantedCourse].sections[wantedSection] push new file
		// filter
		filter := bson.D{{"courses.sections._id", sectionObjId}}
		rs := collection.FindOneAndUpdate(c.Request.Context(), filter, bson.D{{"$push", bson.D{{"courses.sections.$.files", fileObject}}}})
		err = rs.Err()
		if err != nil {
			log.Printf("Error inserting file: %v", err)
			err = session.AbortTransaction(c.Request.Context())
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNABLE_CREATE_FILE})
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
func deleteFileObject(c *gin.Context, collection *mongo.Collection, file models.Files) error {
	filter := bson.D{{"courses.sections.files._id", file.ID}}
	rs := collection.FindOneAndUpdate(c.Request.Context(), filter, bson.D{{"$pull", bson.D{{"courses.sections.$.files", bson.D{{"_id", file.ID}}}}}})
	err := rs.Err()
	if err != nil {
		log.Printf("Error deleting file object: %v", err)
		return err

	}
	return nil
}

// For the file (the actual file, not the file object in the database), we need to implement the following functions:
func deleteSavedFile(filename string) error {
	dir, err := GetStorageFile("files")
	if err != nil {
		log.Printf("Error getting storage file: %v", err)
	}
	err = os.Remove(path.Join(dir, filename))
	if err != nil {
		log.Printf("Error deleting file: %v", err.Error())
	}
	return err
}

func DeleteFileTransaction(client *mongo.Client, collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {

		var file models.Files
		fileId := c.Param("id")
		fileObjectId, errD := primitive.ObjectIDFromHex(fileId)
		if errD != nil {
			log.Printf("Error converting file id: %v", errD)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		file.ID = fileObjectId
		session, err := client.StartSession()
		if err != nil {
			log.Printf("Error starting session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return

		}
		defer func() {
			session.EndSession(c.Request.Context())
			return
		}()
		errDF := deleteSavedFile(file.Url)

		if errDF != nil {
			log.Printf("Error deleting file object: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return
		}
		errOF := deleteFileObject(c, collection, file)
		if errOF != nil {
			log.Printf("Error deleting file object: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return

		}
		err = session.CommitTransaction(c.Request.Context())
		if err != nil {
			log.Printf("Error deleting saved file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": shared.FILE_DELETED})
		return
	}
}
