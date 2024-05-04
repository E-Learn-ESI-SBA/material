package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
	"os"
	"path"
	"time"
)

func EditFile(ctx context.Context, collection *mongo.Collection, file models.Files) error {
	up, err := collection.UpdateOne(ctx, bson.M{"_id": file.ID}, bson.M{"$set": bson.M{"name": file.Name}})
	if err != nil {
		log.Println("Error updating file: ", err)
		return errors.New(shared.FILE_NOT_UPDATED)
	}
	if up.ModifiedCount == 0 {
		log.Println("File Note Found ", err)
		return errors.New(shared.FILE_NOT_FOUND)
	}

	return nil
}

func CreateFileObject(ctx context.Context, collection *mongo.Collection, sectionId primitive.ObjectID, file models.Files) error {
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

func DeleteFileObject(ctx context.Context, collection *mongo.Collection, fileId primitive.ObjectID) error {

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

type FileResponse struct {
	models.Module
	File models.Files `bson:"file"`
}

func GetFileObject(ctx context.Context, collection *mongo.Collection, fileId primitive.ObjectID) (FileResponse, error) {

	var module FileResponse
	pipeline := bson.A{
		bson.M{"$unwind": "$courses"},
		bson.M{"$unwind": "$courses.sections"},
		bson.M{"$unwind": "$courses.sections.files"},
		bson.M{"$match": bson.M{"courses.sections.files._id": fileId}},
		bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": []interface{}{"$$ROOT", bson.M{"file": bson.M{"_id": "$courses.sections.files._id", "url": "$courses.sections.files.url"}}}}}},
		bson.M{"$project": bson.M{"courses": 0}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return module, err

	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err := cursor.Decode(&module)
		if err != nil {
			return module, err
		}

	}
	return module, nil

}

func DeleteSavedFile(filename string, dir string) error {
	filePath := path.Join(dir, filename)
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Error deleting file: %v", err.Error())
	}
	return err
}
