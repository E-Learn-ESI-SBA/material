package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/models"
)

type ExtendedSection struct {
	models.Section

	Files    []models.Files   `json:"files"`
	Videos   []models.Video   `json:"videos"`
	Lectures []models.Lecture `json:"contents"`
}

func GetSectionDetailsById(ctx context.Context, collection *mongo.Collection, sectionId string) (ExtendedSection, error) {
	var sections ExtendedSection
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"_id": sectionId},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "videos",
				"localField":   "_id",
				"foreignField": "section_id",
				"as":           "videos",
			},
		}, bson.M{
			"$lookup": bson.M{
				"from":         "files",
				"localField":   "_id",
				"foreignField": "section_id",
				"as":           "files",
			}}, bson.M{
			"$lookup": bson.M{
				"from":         "lectures",
				"localField":   "_id",
				"foreignField": "section_id",
				"as":           "lectures",
			},
		},
	}
	cursor, err := collection.Find(ctx, pipeline)
	if err != nil {
		log.Printf("Error While Getting Lectures By Module: %v\n", err)
		return sections, err

	}
	cursorError := cursor.All(ctx, &sections)
	if cursorError != nil {
		log.Printf("Error While Parsing Lectures By Module: %v\n", cursorError)
		return sections, cursorError
	}
	return sections, nil
}

type ExtendedSectionWithNote struct {
	ExtendedSection
	Notes []models.Note `json:"notes"`
}

func GetSectionFromStudent(ctx context.Context, SectionCollection *mongo.Collection, noteCollection *mongo.Collection, sectionId string, studentId string) (models.Section, error) {
	// get notes and extend the section in parallel

	lightSection, err := GetSectionDetailsById(ctx, SectionCollection, sectionId)
	if err != nil {
		log.Printf("Error While Getting Section By Student: %v\n", err)
		return lightSection.Section, err
	}

	sections := ExtendedSectionWithNote{
		ExtendedSection: lightSection,
	}

}
