package transactions

import (
	"context"
	"errors"
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
		err = createFileObject(c.Request.Context(), collection, sectionObjId, fileObject)
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
func createFileObject(ctx context.Context, collection *mongo.Collection, sectionId primitive.ObjectID, file models.Files) error {
	file.CreatedAt = time.Now()
	file.UpdatedAt = file.CreatedAt

	filter := bson.D{{"courses.sections._id", sectionId}}
	update := bson.M{
		"$push": bson.M{
			"courses.$[course].sections.$[section].files": file,
		},
	}
	arrayFilters := bson.A{
		bson.M{"course.sections._id": sectionId},
		bson.M{"section._id": sectionId},
	}
	opts := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
	rs := collection.FindOneAndUpdate(ctx, filter, update, opts)
	err := rs.Err()
	if err != nil {
		return errors.New(shared.UNABLE_CREATE_FILE)
	}
	return nil
}
func deleteFileObject(ctx context.Context, collection *mongo.Collection, fileId primitive.ObjectID) error {

	filter := bson.D{{"courses.sections.files._id", fileId}}
	update := bson.M{
		"$pull": bson.M{
			"courses.$[course].sections.$[section].files": bson.M{"_id": fileId},
		},
	}
	arrayFilters := bson.A{
		bson.M{"course.sections.files._id": fileId},
		bson.M{"section.files._id": fileId},
	}
	opts := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
	rs := collection.FindOneAndUpdate(ctx, filter, update, opts)
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
		rs := collection.FindOne(ctx, bson.D{{"courses.sections.files._id", fileObjectId}})
		err := rs.Err()
		if err != nil {
			log.Printf("File not found error : %v", err.Error())
			c.JSON(http.StatusNotFound, shared.FILE_NOT_FOUND)
		}
		rs.Decode(&file)
		session, errS := client.StartSession()
		if errS != nil {
			log.Printf("Error starting session: %v", errS)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return

		}
		defer func() {
			session.EndSession(ctx)
			return
		}()
		errDF := deleteSavedFile(file.Url)

		if errDF != nil {
			log.Printf("Error deleting file object: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_DELETED})
			return
		}
		errOF := deleteFileObject(c.Request.Context(), collection, fileObjectId)
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
