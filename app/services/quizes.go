package services

import (
	"context"
	"log"
	"madaurus/dev/material/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



func CreateQuiz(ctx context.Context, collection *mongo.Collection, quiz models.Quiz) error {
	quiz.ID = primitive.NewObjectID() // necessary with mongo docker image
	_, err := collection.InsertOne(ctx, quiz)
	if err != nil {
		log.Printf("Error While Creating Quiz: %v\n", err)
		return err
	}
	return nil
}


func UpdateQuiz(ctx context.Context, collection *mongo.Collection, quiz models.Quiz, teacherID int) error {
	_, err := collection.UpdateOne(ctx, bson.M{"_id": quiz.ID, "teacher_id": teacherID}, bson.D{{"$set", quiz}})
	if err != nil {
		log.Printf("Error While Updating Quiz: %v\n", err)
		return err
	}
	return nil
}


func DeleteQuiz(ctx context.Context, collection *mongo.Collection, quizID string, teacherID int) error {
	a, err := collection.DeleteOne(ctx, bson.D{{"_id", quizID}, {"teacher_id", teacherID}})
	log.Println(quizID, teacherID, a.DeletedCount)
	if err != nil {
		log.Printf("Error While Deleting Course: %v\n", err)
		return err
	}
	return nil
}


func GetQuiz(ctx context.Context, collection *mongo.Collection, quizID string) (models.Quiz, error) {
	var quiz models.Quiz
	objectId, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return quiz, err
	}
	filter := bson.D{{"_id", objectId}}
	err = collection.FindOne(ctx, filter).Decode(&quiz)
	log.Println("fetched quiz", quiz)
	if err != nil {
		log.Printf("Error While Getting Quiz: %v\n", err)
		return quiz, err
	}
	return quiz, nil
}


func GetQuizesByAdmin(ctx context.Context, collection *mongo.Collection) ([]models.Quiz, error) {
	var quizes []models.Quiz
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error While Getting Quizes: %v\n", err)
		return quizes, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var quiz models.Quiz
		cursor.Decode(&quiz)
		quizes = append(quizes, quiz)
	}
	return quizes, nil
}