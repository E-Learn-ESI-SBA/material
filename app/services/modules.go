package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
)

func GetSectionsByTeacher(ctx context.Context, collection *mongo.Collection, teacherId string) ([]models.Section, error) {
	var courses []models.Section
	cursor, err := collection.Find(ctx, bson.D{{"teacher_id", teacherId}})
	if err != nil {
		return nil, err
	}
	cursorError := cursor.All(ctx, &courses)
	if cursorError != nil {
		return nil, cursorError
	}
	return courses, nil
}

func GetSectionsByAdmin(ctx context.Context, collection *mongo.Collection) ([]models.Section, error) {
	var courses []models.Section
	cursor, err := collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	cursorError := cursor.All(ctx, &courses)
	if cursorError != nil {
		return nil, cursorError
	}
	return courses, nil
}

// GetModulesByFilter Basic Usage  : GetModulesByFilter(ctx, collection, filterStruct, "public", nil) for public endpoints
// Advanced Usage: GetModulesByFilter(ctx, collection, filterStruct, "private", &teacherId) for private endpoints
func GetModulesByFilter(ctx context.Context, collection *mongo.Collection, filterStruct interfaces.ModuleFilter, usage string, teacherId *string) ([]models.Module, error) {
	var modules []models.Module
	var filter bson.D
	if usage == "public" {
		filter = bson.D{{"year", filterStruct.Year}, {"semester", filterStruct.Semester}, {"speciality", filterStruct.Speciality}, {
			"isPublic", true}}

	} else if teacherId != nil {
		filter = bson.D{{"year", filterStruct.Year}, {"semester", filterStruct.Semester}, {"speciality", filterStruct.Speciality}, {
			"teacher_id", *teacherId}}
	} else {
		return nil, errors.New("teacher Id is required for this operation")
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	cursorError := cursor.All(ctx, &modules)
	if cursorError != nil {
		return nil, cursorError
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("failed to close cursor")
		}
	}(cursor, ctx)
	return modules, nil
}

func EditModuleVisibility(ctx context.Context, collection *mongo.Collection, moduleId string, visibility bool) error {
	filter := bson.D{{"_id", moduleId}}
	update := bson.D{{"$set", bson.D{{"isPublic", visibility}}}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
