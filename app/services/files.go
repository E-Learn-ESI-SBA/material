package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
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
