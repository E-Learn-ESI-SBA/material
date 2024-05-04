package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
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

func GetSectionDetailsById(ctx context.Context, collection *mongo.Collection, sectionId string, pips ...bson.M) (models.ExtendedSection, error) {
	var sections models.ExtendedSection
	id, errId := primitive.ObjectIDFromHex(sectionId)
	if errId != nil {
		log.Printf("Error While Parsing Section ID: %v\n", errId)
		return sections, errId
	}
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"_id": id},
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
	cursor, err := collection.Aggregate(ctx, pipeline)
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
func GetSectionFromStudent(ctx context.Context, SectionCollection *mongo.Collection, sectionId string, studentId string) (models.ExtendedSection, error) {
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
		return models.ExtendedSection{}, err
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
func EditSection(ctx context.Context, collection *mongo.Collection, section models.Section, sectionId primitive.ObjectID, teacherId string) error {
	// update only the name
	rs := collection.FindOneAndUpdate(ctx, bson.D{{"courses.sections._id", sectionId}, {"teacher_id", teacherId}}, bson.D{{"$set", bson.D{{"courses.sections.$.name", section.Name}}}})
	err := rs.Err()
	if err != nil {
		log.Printf("Error updating section: %v", err)
		return errors.New(shared.UNABLE_UPDATE_SECTION)
	}
	return nil
}

func CreateSection(ctx context.Context, collection *mongo.Collection, section models.Section, courseId primitive.ObjectID) error {
	section.Files = []models.Files{}
	section.Videos = []models.Video{}
	section.Lectures = []models.Lecture{}
	rs := collection.FindOneAndUpdate(ctx, bson.D{{"courses._id", courseId}}, bson.D{{"$push", bson.D{{"courses.$.sections", section}}}})
	err := rs.Err()
	if err != nil {
		log.Printf("Error inserting section: %v", err)
		return errors.New(shared.UNABLE_CREATE_SECTION)
	}
	return nil
}

func DeleteSection(ctx context.Context, collection *mongo.Collection, sectionId primitive.ObjectID, teacherId string) error {
	rs := collection.FindOneAndUpdate(ctx, bson.D{{"teacher_id", teacherId}, {"courses.sections._id", sectionId}, {
		"courses.sections.files", bson.D{{"$size", 0}},
	},
		{"courses.sections.videos", bson.D{{"$size", 0}}},
		{"courses.sections.lectures", bson.D{{"$size", 0}}},
	}, bson.D{{"$pull", bson.D{{"courses.sections", bson.D{{"_id", sectionId}}}}}})
	err := rs.Err()
	if err != nil {
		log.Printf("Error While Deleting Section: %v\n", err)
		return errors.New(shared.UNABLE_DELETE_SECTION)

	}
	return nil
}

/*
func DeleteSection(ctx context.Context, collection *mongo.Collection, sectionId primitive.ObjectID) error {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"_id": bson.M{"$eq": sectionId}},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "files",
				"localField":   "_id",
				"foreignField": "sectionID",
				"as":           "files",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "videos",
				"localField":   "_id",
				"foreignField": "sectionID",
				"as":           "videos",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "lectures",
				"localField":   "_id",
				"foreignField": "sectionID",
				"as":           "lectures",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":         "$files",
				"preserveNull": true,
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":         "$videos",
				"preserveNull": true,
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":         "$lectures",
				"preserveNull": true,
			},
		},
		bson.M{
			"$project": bson.M{
				"_id": 0, // exclude original document
				"count": bson.M{
					"$sum": bson.A{1, "$files", "$videos", "$lectures"},
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$count"},
			},
		},
	}
	var count int64 = -1
	result, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error While Trying To Deeply delete Section %v", err.Error())

	}

	errR := result.Decode(&count)
	if errR != nil {
		log.Printf("Error While Trying To Deeply delete Section %v", errR.Error())
	}

	if count > 0 {
		return errors.New(shared.UNABLE_DELETE_SECTION)
	}
	d, err := collection.DeleteOne(ctx, bson.D{{"_id", sectionId}})

	if err != nil || d.DeletedCount < 1 {
		return errors.New(shared.UNABLE_DELETE_SECTION)

	}
	return nil
}


*/
