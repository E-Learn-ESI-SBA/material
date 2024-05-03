package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
)

func CreateCourse(ctx context.Context, collection *mongo.Collection, course models.Course, moduleId primitive.ObjectID) error {
	// instead of insert , find the module from moduleId  , then insert the course in the courses array of the module
	rs, err := collection.UpdateOne(ctx, bson.D{{"_id", moduleId}}, bson.D{{"$push", bson.D{{"courses", course}}}})
	if err != nil {
		log.Printf("Error While Creating Course: %v\n", err)
		return err

	}
	if rs.ModifiedCount == 0 {
		log.Printf("Error While Creating Course: %v\n", err)
		return errors.New(shared.UNABLE_CREATE_COURSE)
	}

	/*
		result, err := collection.InsertOne(ctx, course)
		log.Printf("Course Created: %v\n", result.InsertedID)
		if err != nil {
			log.Printf("Error While Creating Course: %v\n", err)
			return err
		}

	*/
	return nil
}
func UpdateCourse(ctx context.Context, collection *mongo.Collection, course models.Course, teacherId string, moduleId primitive.ObjectID) error {

	// search first collection module by module Id, then course from course array by course.id , then update the course
	rs, err := collection.UpdateOne(ctx, bson.D{{"_id", moduleId}, {"courses._id", course.ID}}, bson.D{{"$set", bson.D{{"courses.$", course}}}})
	if err != nil {
		log.Printf("Error While Updating Course: %v\n", err)
		return errors.New(shared.UNABLE_UPDATE_COURSE)
	}
	if rs.ModifiedCount == 0 {
		return errors.New(shared.UNABLE_UPDATE_COURSE)
	}
	return nil

}

// GetCoursesByInstructor is a function that returns a list of courses that an instructor is teaching
func GetCoursesByInstructor(ctx context.Context, collection *mongo.Collection, instructorID string) ([]models.Module, error) {
	var modules []models.Module
	cursor, err := collection.Find(ctx, bson.D{{"teacher_id", instructorID}})
	if err != nil {
		log.Printf("Error While Getting Courses By Instructor: %v\n", err)
		return nil, errors.New(shared.UNABLE_GET_COURSES)

	}
	cursor.Decode(&modules)
	return modules, nil
}

/*
	func DeleteCourse(ctx context.Context, collection *mongo.Collection, courseID string, teacherId int) error {
		_, err := collection.DeleteOne(ctx, bson.D{{"_id", courseID}, {"teacher_id", teacherId}})
		if err != nil {
			log.Printf("Error While Deleting Course: %v\n", err)
			return err
		}
		return nil
	}
*/

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

func DeleteCourse(ctx context.Context, collection *mongo.Collection, courseId primitive.ObjectID, moduleId primitive.ObjectID) error {
	rs, err := collection.UpdateOne(ctx, bson.D{{"_id", moduleId}, {"courses._id", courseId}, {
		"courses.sections", bson.D{{"$size", 0}}},
	}, bson.D{{"$pull", bson.D{{"courses", bson.D{{"_id", courseId}}}}}})
	if err != nil {
		log.Printf("Error While Deleting Course: %v\n", err)
		return errors.New(shared.UNABLE_DELETE_COURSE)

	}
	if rs.ModifiedCount == 0 {
		return errors.New(shared.UNABLE_DELETE_COURSE)

	}
	return nil
}
func GetCourseById(ctx context.Context, collection *mongo.Collection, courseId primitive.ObjectID) (models.Course, error) {
	var course models.Course
	err := collection.FindOne(ctx, bson.D{{"courses._id", courseId}}).Decode(&course)
	if err != nil {
		log.Printf("Error While Getting Course By ID: %v\n", err)
		return models.Course{}, errors.New(shared.UNABLE_GET_COURSE)
	}
	return course, nil
}

/*

pip := bson.A{
		bson.M{
			"$match": bson.M{"_id": courseId},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "sections",
				"localField":   "_id",
				"foreignField": "course_id",
				"as":           "sections",
			},
		},
		bson.M{
			"$unwind": "$sections",
		},
		bson.M{
			"$project": bson.M{
				"_id": 0, // exclude original document
				"count": bson.M{
					"$sum": bson.A{1, "$sections"},
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
	cursor, err := collection.Aggregate(ctx, pip)
	if err != nil {
		log.Printf("Error While Deleting Course: %v\n", err)
		return errors.New(shared.UNABLE_DELETE_COURSE)
	}
	cursorError := cursor.Decode(&count)
	if cursorError != nil {
		log.Printf("Error While Parsing Course Count: %v\n", cursorError)
		return errors.New(shared.UNABLE_DELETE_COURSE)
	}
	if count > 0 {
		return errors.New(shared.UNABLE_DELETE_COURSE)
	}
	_, err = collection.DeleteOne(ctx, bson.D{{"_id", courseId}})
	if err != nil {
		log.Printf("Error While Deleting Course: %v\n", err)
		return errors.New(shared.UNABLE_DELETE_COURSE)
	}
	return nil
*/
