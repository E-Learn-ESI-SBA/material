package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"madaurus/dev/material/app/models"
)

func GetSectionsByCourse(ctx context.Context, collection *mongo.Collection, moduleId string) ([]models.Section, error) {
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

type ExtendedSection struct {
	models.Section

	Files    []models.Files        `json:"files"`
	Videos   []models.Video        `json:"videos"`
	Lectures []models.Lecture      `json:"contents"`
	Notes    *[]models.StudentNote `json:"note"`
}

func GetSectionDetailsById(ctx context.Context, collection *mongo.Collection, sectionId string, pips ...bson.M) (ExtendedSection, error) {
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
	for _, pip := range pips {
		pipeline = append(pipeline, pip)
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

func GetSectionFromStudent(ctx context.Context, SectionCollection *mongo.Collection, sectionId string, studentId int) (ExtendedSection, error) {
	pip := bson.M{
		"$lookup": bson.M{
			"from":         "student_notes",
			"localField":   "_id",
			"foreignField": "section_id",
			"as":           "notes",
		},
	}
	extendedSection, err := GetSectionDetailsById(ctx, SectionCollection, sectionId, pip)
	if err != nil {
		log.Printf("Error While Getting Section By Student: %v\n", err)
		return ExtendedSection{}, err
	}
	filteredNotes := make([]models.StudentNote, 0)
	notes := extendedSection.Notes
	for _, note := range *notes {
		if note.StudentID == studentId {
			filteredNotes = append(filteredNotes, note)
		}
	}
	extendedSection.Notes = &filteredNotes
	return extendedSection, nil
}

func EditSection(ctx context.Context, collection *mongo.Collection, section models.Section, sectionId string) error {
	filter := bson.D{{"_id", sectionId}}
	update := bson.D{{"$set", section}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error While Updating Section: %v\n", err)
		return err
	}
	return nil
}

func CreateSection(ctx context.Context, collection *mongo.Collection, section models.Section) error {
	_, err := collection.InsertOne(ctx, section)
	if err != nil {
		log.Printf("Error While Creating Section: %v\n", err)
		return err
	}
	return nil
}

func DeleteSection(ctx context.Context, collection *mongo.Collection, sectionId string) error {
	filter := bson.D{{"_id", sectionId}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error While Deleting Section: %v\n", err)
		return err
	}
	return nil
}
