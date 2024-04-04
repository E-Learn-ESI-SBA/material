package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
)

// GetModulesByFilter Basic Usage  : GetModulesByFilter(ctx, collection, filterStruct, "public", nil) for public endpoints
// Advanced Usage: GetModulesByFilter(ctx, collection, filterStruct, "private", &teacherId) for private endpoints
func GetModulesByFilter(ctx context.Context, collection *mongo.Collection, filterStruct interfaces.ModuleFilter, usage string, teacherId *int) ([]models.Module, error) {
	var modules []models.Module
	var filter bson.D
	opts := options.Find().SetProjection(bson.D{{"courses", 0}})
	if usage == "public" {
		filter = bson.D{{"year", filterStruct.Year}, {"semester", filterStruct.Semester}, {"speciality", filterStruct.Speciality}, {
			"isPublic", true}}

	} else if teacherId != nil {
		filter = bson.D{{"year", filterStruct.Year}, {"semester", filterStruct.Semester}, {"speciality", filterStruct.Speciality}, {
			"teacher_id", *teacherId}}
	} else {
		return nil, errors.New("teacher Id is required for this operation")
	}
	cursor, err := collection.Find(ctx, filter, opts)
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

func DeleteModule(ctx context.Context, collection *mongo.Collection, moduleId string, teacherId int) error {
	filter := bson.D{{"_id", moduleId}, {"teacher_id", teacherId}}
	err := collection.FindOneAndDelete(ctx, filter)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return errors.New("error while trying to delete the module")

	}
	return nil
}

func UpdateModule(ctx context.Context, collection *mongo.Collection, module models.Module) error {
	filter := bson.D{{"_id", module.ID}}
	// do not update the courses
	update := bson.D{{"$set", module}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return errors.New("error while trying to update the module")
	}
	return nil
}

func CreateModule(ctx context.Context,collection *mongo.Collection, module models.Module) error {
	_, err := collection.InsertOne(ctx, module)
	if err != nil {
		log.Printf("error while trying to create the module")
	}
	return err
}

func GetModuleById(ctx context.Context, collection *mongo.Collection, moduleId string) (models.ExtendedModule, error) {
	// make aggregation to get the courses
	// then select sections from the courses
	// then select the lectures from the sections and videos from sections
	var module models.ExtendedModule

	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"_id": moduleId},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "courses",
				"localField":   "_id",
				"foreignField": "module_id",
				"as":           "courses",
			},
		},
		bson.M{
			"$unwind": "$courses",
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "sections",
				"localField":   "courses._id",
				"foreignField": "course_id",
				"as":           "courses.sections",
			},
		},
		bson.M{
			"$unwind": "$courses.sections",
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "lectures",
				"localField":   "courses.sections._id",
				"foreignField": "section_id",
				"as":           "courses.sections.lectures",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "videos",
				"localField":   "courses.sections._id",
				"foreignField": "section_id",
				"as":           "courses.sections.videos",
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error While Getting module details")
		return models.ExtendedModule{}, err

	}
	errCursor := cursor.All(ctx, &module)
	if errCursor != nil {
		return models.ExtendedModule{}, errCursor

	}
	return module, nil
}
