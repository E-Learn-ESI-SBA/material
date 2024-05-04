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
	"time"
)

func GetTeacherLecture(collection *mongo.Collection, ctx context.Context, lectureId primitive.ObjectID) (models.Lecture, error) {
	var lecture models.Lecture
	filter := bson.D{{"courses.sections.lectures._id", lectureId}}
	// from module with  courses.sections.lectures , select only the lecture with the id lectureId
	opts := options.FindOne().SetProjection(bson.D{{"courses.sections.lectures.$", 1}, {"_id", 1}})
	err := collection.FindOne(ctx, filter, opts).Decode(&lecture)
	if err != nil {
		log.Printf("Error While Getting Lecture: %v\n", err)
		return models.Lecture{}, err
	}
	return lecture, nil
}

func CreateLecture(collection *mongo.Collection, ctx context.Context, lecture models.Lecture, sectionId primitive.ObjectID) error {
	lecture.ID = primitive.NewObjectID()
	lecture.CreatedAt = time.Now()
	log.Printf("Section id: %v\n", sectionId)
	lecture.UpdatedAt = lecture.CreatedAt
	filter := bson.M{"courses.sections._id": sectionId}
	update := bson.M{
		"$push": bson.M{
			"courses.$[course].sections.$[section].lectures": lecture,
		},
	}
	arrayFilters := bson.A{
		bson.M{"course.sections._id": sectionId},
		bson.M{"section._id": sectionId},
	}
	opts := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
	rs := collection.FindOneAndUpdate(ctx, filter, update, opts)
	err := rs.Err()
	if err != nil {
		log.Printf("Error While Creating Lecture: %v\n", err)
		return errors.New(shared.LECTURE_NOT_CREATED)
	}
	return nil
}

func UpdateLecture(collection *mongo.Collection, ctx context.Context, lecture models.Lecture) error {
	filter := bson.D{{"courses.sections.lectures._id", lecture.ID}}
	update := bson.D{{"$set", bson.D{{"courses.sections.lectures.$", lecture}}}}
	rs, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error While Updating Lecture: %v\n", err)
		return errors.New(shared.LECTURE_NOT_UPDATED)
	}
	if rs.ModifiedCount == 0 {
		log.Printf("Error While Updating Lecture: \n")
		return errors.New(shared.LECTURE_NOT_UPDATED)

	}
	return nil
}

func DeleteLecture(collection *mongo.Collection, ctx context.Context, lectureId primitive.ObjectID) error {
	filter := bson.D{{"courses.sections.lectures._id", lectureId}}
	update := bson.D{{"$pull", bson.D{{"courses.sections.lectures", bson.D{{"_id", lectureId}}}}}}
	rs, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error While Deleting Lecture: %v\n", err)
		return errors.New(shared.LECTURE_NOT_DELETED)
	}
	if rs.ModifiedCount == 0 {
		log.Printf("Error While Deleting Lecture: \n")
		return errors.New(shared.LECTURE_NOT_DELETED)

	}
	return nil
}
