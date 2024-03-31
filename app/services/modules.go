package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/models"
)

func CreateModule(ctx context.Context, collection *mongo.Collection, module models.Module) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(ctx, module)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func updateModule(ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func DeleteModuleById(ctx context.Context, collection *mongo.Collection, moduleId string, teacherId string) error {
	filter := bson.D{{"teacher_id", teacherId}, {"_id", moduleId}}
	var modules models.Module
	result := collection.FindOneAndDelete(ctx, filter).Decode(&modules)
	if result != nil {
		log.Printf("Error while deleting module %v", result)
	}
	return result
}
func GetModulesByStudent(ctx context.Context, collection *mongo.Collection, filter bson.D) ([]models.Module, error) {
	var modules []models.Module
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	cursorError := cursor.All(ctx, &modules)
	if cursorError != nil {
		return nil, cursorError
	}
	return modules, nil
}
