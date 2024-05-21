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
	"time"
)

type ResponseLecture struct {
	models.Module
	lecture models.Lecture `bson:"lecture"`
}

func GetTeacherLecture(collection *mongo.Collection, ctx context.Context, lectureId primitive.ObjectID) (models.Lecture, error) {
	var lecture ResponseLecture
	pipeline := bson.A{
		bson.M{"$unwind": "$courses"},
		bson.M{"$unwind": "$courses.sections"},
		bson.M{"$unwind": "$courses.sections.lectures"},
		bson.M{"$match": bson.M{"courses.sections.lectures._id": lectureId}},
		bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": []interface{}{"$$ROOT", bson.M{"lecture": bson.M{"_id": "$courses.sections.lectures._id", "group": "$courses.sections.lectures.group", "name": "$courses.sections.lectures.name", "content": "$courses.sections.lectures.content"}}}}}},
		bson.M{"$project": bson.M{"courses": 0}},
	}
	//	filter := bson.D{{"courses.sections.lectures._id", lectureId}}
	// from module with  courses.sections.lectures , select only the lecture with the id lectureId
	//	opts := options.FindOne().SetProjection(bson.D{{"courses.sections.lectures.$", 1}, {"_id", 1}})
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error While Getting Lecture: %v\n", err)
		return models.Lecture{}, err
	}
	for cursor.Next(ctx) {
		err := cursor.Decode(&lecture)
		if err != nil {
			log.Printf("Error While Decoding Lecture: %v\n", err)
			return models.Lecture{}, err
		}

	}
	return lecture.lecture, nil
}

func CreateLecture(collection *mongo.Collection, ctx context.Context, lecture models.Lecture, sectionId primitive.ObjectID, permitApi *permit.Client, client *mongo.Client) error {
	session, err := client.StartSession()
	if err != nil {
		log.Printf("Error While Creating Lecture: %v\n", err)
		return errors.New(shared.LECTURE_NOT_CREATED)

	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {

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
		err = rs.Err()
		if err != nil {
			sessionContext.AbortTransaction(ctx)
			log.Printf("Error While Creating Lecture: %v\n", err)
			return nil, errors.New(shared.LECTURE_NOT_CREATED)
		}
		sectionIdStr := sectionId.Hex()
		err = utils.CreateResourceInstance(permitApi, "lectures", lecture.ID.Hex(), &sectionIdStr, &iam.SECTIONS, &iam.PARENT)
		if err != nil {
			sessionContext.AbortTransaction(ctx)
			log.Printf("Error While Creating Resource Instance: %v\n", err)
			return nil, errors.New(shared.LECTURE_NOT_CREATED)
		}
		return nil, nil
	})
	return err
}

func UpdateLecture(collection *mongo.Collection, ctx context.Context, lecture models.Lecture) error {
	filter := bson.D{{"courses.sections.lectures._id", lecture.ID}}
	update := bson.M{
		"$set": bson.M{
			"courses.$[course].sections.$[section].lectures.$[lecture]": lecture,
		},
	}
	arrayFilters := bson.A{
		bson.M{"course.sections.lectures._id": lecture.ID},
		bson.M{"section.lectures._id": lecture.ID},
		bson.M{"lecture._id": lecture.ID},
	}
	opts := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
	rs := collection.FindOneAndUpdate(ctx, filter, update, opts)
	err := rs.Err()
	if err != nil {
		log.Printf("Error While Updating Lecture: %v\n", err)
		return errors.New(shared.LECTURE_NOT_UPDATED)
	}

	return nil
}

func DeleteLecture(collection *mongo.Collection, ctx context.Context, lectureId primitive.ObjectID) error {
	filter := bson.D{{"courses.sections.lectures._id", lectureId}}
	update := bson.M{
		"$pull": bson.M{
			"courses.$[course].sections.$[section].lectures.$[lecture]._id": lectureId,
		},
	}
	arrayFilters := bson.A{
		bson.M{"course.sections.lectures._id": lectureId},
		bson.M{"section.lectures._id": lectureId},
		bson.M{"lecture._id": lectureId},
	}
	//update := bson.D{{"$pull", bson.D{{"courses.sections.lectures", bson.D{{"_id", lectureId}}}}}}

	opts := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
	rs := collection.FindOneAndUpdate(ctx, filter, update, opts)
	err := rs.Err()
	if err != nil {
		log.Printf("Error While Deleting Lecture: %v\n", err)
		return errors.New(shared.LECTURE_NOT_DELETED)
	}
	return nil
}
