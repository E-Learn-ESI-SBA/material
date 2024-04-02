package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/models"
)

func GetTeacherLecture(collection *mongo.Collection, ctx context.Context, lectureId string) (models.Lecture, error) {
	var lecture models.Lecture
	filter := bson.D{{"_id", lectureId}}
	err := collection.FindOne(ctx, filter).Decode(&lecture)
	if err != nil {
		log.Printf("Error While Getting Lecture: %v\n", err)
		return models.Lecture{}, err
	}
	return lecture, nil
}

func CreateLecture(collection *mongo.Collection, ctx context.Context, lecture models.Lecture) error {
	_, err := collection.InsertOne(ctx, lecture)
	if err != nil {
		log.Printf("Error While Creating Lecture: %v\n", err)
		return err
	}
	return nil
}

func UpdateLecture(collection *mongo.Collection, ctx context.Context, lecture models.Lecture) error {
	filter := bson.D{{"_id", lecture.ID}}
	update := bson.D{{"$set", lecture}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error While Updating Lecture: %v\n", err)
		return err
	}
	return nil
}

func DeleteLecture(collection *mongo.Collection, ctx context.Context, lectureId string) error {
	filter := bson.D{{"_id", lectureId}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error While Deleting Lecture: %v\n", err)
		return err
	}
	return nil
}
