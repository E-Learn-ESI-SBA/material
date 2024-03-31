package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"madaurus/dev/material/app/models"
)

func GetSectionsByModule(ctx context.Context, collection *mongo.Collection, moduleId string) ([]models.Section, error) {
	var sections []models.Section
	func(usedSection *[]models.Section) {
		filter := bson.D{{"module_id", moduleId}}
		sort := bson.D{{"created_at", 1}}
		opts := options.Find().SetSort(sort)
		cursor, err := collection.Find(ctx, filter, opts)
		if err != nil {
			log.Printf("Error While Getting Sections By Module: %v\n\n", err)
			return
		}
		cursorError := cursor.All(ctx, &usedSection)
		if cursorError != nil {
			log.Printf("Error While Parsing Sections By Module: %v\n\n", cursorError)
			return
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {

			}
		}(cursor, ctx)

	}(&sections)

	return sections, nil
}
