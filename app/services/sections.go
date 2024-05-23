package services

import (
	"context"
	"errors"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/shared/iam"
	"madaurus/dev/material/app/utils"
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

/*
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
*/
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

func CreateSection(ctx context.Context, collection *mongo.Collection, section models.Section, courseId primitive.ObjectID, premitApi *permit.Client, client *mongo.Client) error {
	section.Files = []models.Files{}
	section.Videos = []models.Video{}
	section.Lectures = []models.Lecture{}
	session, err := client.StartSession()
	if err != nil {
		log.Printf("Error starting session: %v", err)
		return errors.New(shared.UNABLE_CREATE_SECTION)

	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		rs := collection.FindOneAndUpdate(ctx, bson.D{{"courses._id", courseId}}, bson.D{{"$push", bson.D{{"courses.$.sections", section}}}})
		err = rs.Err()
		if err != nil {
			sessionCtx.AbortTransaction(ctx)
			log.Printf("Error inserting section: %v", err)
			return nil, errors.New(shared.UNABLE_CREATE_SECTION)
		}
		courseIdStr := courseId.Hex()
		err = utils.CreateResourceInstance(premitApi, "sections", section.ID.Hex(), &courseIdStr, &iam.CHAPTERS, &iam.PARENT)
		if err != nil {
			sessionCtx.AbortTransaction(ctx)
			log.Printf("Error creating resource instance: %v", err)
			return nil, errors.New(shared.UNABLE_CREATE_SECTION)
		}
		return nil, nil
	})

	return err
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

type SectionWithChapterName struct {
	ChapterName string `json:"chapter_name"`
	models.Section
}

func GetSectionsByAdmin(ctx context.Context, collection *mongo.Collection) ([]SectionWithChapterName, error) {
	// Select only sections , and please write  a valid query
	var modules []models.Module
	var sections []SectionWithChapterName

	cursor, err := collection.Find(ctx, bson.D{{}}, options.Find().SetProjection(bson.D{{"courses.sections", 1}}))
	if err != nil {
		log.Printf("Error While Getting Sections By Admin: %v\n", err)
		return sections, err
	}
	cursorError := cursor.All(ctx, &modules)
	if cursorError != nil {
		log.Printf("Error While Parsing Sections By Admin: %v\n", cursorError)
		return sections, cursorError
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("failed to close cursor")
		}
	}(cursor, ctx)
	// Transform it to sections array
	for _, module := range modules {
		for _, course := range module.Courses {
			for _, section := range course.Sections {
				sections = append(sections, SectionWithChapterName{ChapterName: course.Name, Section: section})
			}
		}
	}
	return sections, nil
}
