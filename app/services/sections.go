package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/models"
)

func GetSectionsByCourse(ctx context.Context, collection *mongo.Collection, courseId string) ([]models.Section, error) {
	var sections []models.Section
	filter := bson.D{{"course_id", courseId}}
	sort := bson.D{{}}
	cursor, err := collection.Find(ctx)

}
