package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/models"
)

func CreateCourse(ctx context.Context, collection *mongo.Collection, course models.Course) error {
	_, err := collection.InsertOne(ctx, course)
	if err != nil {
		log.Printf("Error While Creating Course: %v\n", err)
		return err
	}
	return nil
}
func UpdateCourse(ctx context.Context, collection *mongo.Collection, course models.Course, teacherId int) error {
	_, err := collection.UpdateOne(ctx, bson.D{{"_id", course.ID}, {"teacher_id", teacherId}}, bson.D{{"$set", course}})
	if err != nil {
		log.Printf("Error While Updating Course: %v\n", err)
		return err
	}
	return nil
}

// GetCoursesByInstructor is a function that returns a list of courses that an instructor is teaching
func GetCoursesByInstructor(ctx context.Context, collection *mongo.Collection, instructorID int) ([]models.Course, error) {
	// Logic to get courses by instructor
	var courses []models.Course
	cursor, err := collection.Find(ctx, bson.D{{"instructor_id", instructorID}})
	if err != nil {
		log.Printf("Error While Getting Courses By Instructor: %v\n", err)
		return nil, err
	}
	cursorError := cursor.All(ctx, &courses)
	if cursorError != nil {
		log.Printf("Error While Parsing Courses By Instructor: %v\n", cursorError)
		return nil, cursorError
	}
	return courses, nil
}

func DeleteCourse(ctx context.Context, collection *mongo.Collection, courseID string, teacherId int) error {
	_, err := collection.DeleteOne(ctx, bson.D{{"_id", courseID}, {"teacher_id", teacherId}})
	if err != nil {
		log.Printf("Error While Deleting Course: %v\n", err)
		return err
	}
	return nil
}

type ExtendedCourse struct {
	models.Course
	Sections []models.Section `json:"sections"`
}

func FetchCourseSectionsByModule(ctx context.Context, collection mongo.Collection, moduleID string) ([]models.UltraCourse, error) {
	var extendedSections []models.UltraCourse
	id, errId := primitive.ObjectIDFromHex(moduleID)
	if errId != nil {
		log.Printf("Error While Parsing Section ID: %v\n", errId)
		return extendedSections, errId
	}
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"module_id": id},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "sections",
				"localField":   "_id", //
				"foreignField": "course_id",
				"as":           "sections",
			},
		},
		bson.M{
			"$unwind": "$sections", // Unwind sections array to process each section
		},
		// Populate sections with details from Videos, Lectures, Files collections
		bson.M{
			"$lookup": bson.M{
				"from":         "videos",
				"localField":   "sections._id",
				"foreignField": "section_id",
				"as":           "sections.videos",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "lectures",
				"localField":   "sections._id",
				"foreignField": "section_id",
				"as":           "sections.lectures",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "files",
				"localField":   "sections._id",
				"foreignField": "section_id",
				"as":           "sections.files",
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":         "$courses._id",
				"name":        "$courses.name",
				"description": "$courses.description",
				"sections": bson.A{
					bson.M{
						"_id":        "$sections._id",
						"name":       "$sections.name",
						"order":      "$sections.order",
						"teacher_id": "$sections.teacher_id",
						"course_id":  "$sections.course_id",
						"videos":     "$sections.videos",
						"lectures":   "$sections.lectures",
						"files":      "$sections.files",
					},
				},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error While Getting Sections By Module: %v\n", err)
		return nil, err
	}
	cursorError := cursor.All(ctx, &extendedSections)
	if cursorError != nil {
		log.Printf("Error While Parsing Sections By Module: %v\n", cursorError)
		return nil, cursorError
	}
	return extendedSections, nil

}

/*
func GetCoursesByModule(ctx context.Context, collectionCourse mongo.Collection, collectionSection mongo.Collection, moduleID string) (ExtendedCourse, error) {
	var courses []models.Course
	filter := bson.D{{"module_id", moduleID}}
	cursor, err := collectionCourse.Find(ctx, filter)
	if err != nil {
		log.Printf("Error While Getting Course By Module: %v\n", err)
		return ExtendedCourse{}, err
	}
	cursorError := cursor.All(ctx, &courses)
	if cursorError != nil {
		log.Printf("Error While Parsing Course By Module: %v\n", cursorError)
		return ExtendedCourse{}, cursorError

	}
	populateContext, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var course models.Course
	var extendedCourse = ExtendedCourse{Course: course}
	var section models.Section
	for cursor.Next(populateContext) {
		err := cursor.Decode(&course)
		if err != nil {
			log.Printf("Error While Decoding Course By Module: %v\n", err)
			return extendedCourse, err
		}
		sectionCursor, err := collectionSection.Find(populateContext, bson.D{{"course_id", course.ID}})
		if err != nil {
			log.Printf("Error While Getting Sections By Course: %v\n", err)
			return extendedCourse, err
		}
		sectionCursorError := sectionCursor.All(populateContext, &section)
		if sectionCursorError != nil {
			log.Printf("Error While Parsing Sections By Course: %v\n", sectionCursorError)
			return extendedCourse, sectionCursorError

		}
		extendedCourse.Sections = append(extendedCourse.Sections, section)
		err := sectionCursor.Close(populateContext)
		if err != nil {
			return ExtendedCourse{}, err
		}

	}
	defer func() {
		err := cursor.Close(populateContext)
		if err != nil {
			log.Println("failed to close cursor")
		}
	}()
	return extendedCourse, nil

}
*/
