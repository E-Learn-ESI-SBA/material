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
	"madaurus/dev/material/app/utils"
	"os"
	"path"
	"time"
)

type VideoQueryResponse struct {
	models.Module
	Video models.Video `bson:"video"`
}

func CreateVideo(ctx context.Context, collection *mongo.Collection, sectionId primitive.ObjectID, video models.Video) error {

	video.CreatedAt = time.Now()
	video.UpdatedAt = video.CreatedAt
	filter := bson.D{{"courses.sections._id", sectionId}}
	update := bson.M{
		"$push": bson.M{
			"courses.$[course].sections.$[section].videos": video,
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
		return errors.New(shared.UNABLE_TO_CREATE_VIDEO)
	}
	return nil
}
func GetVideo(ctx context.Context, collection *mongo.Collection, videoId primitive.ObjectID) (models.Video, error) {
	var videoResponse VideoQueryResponse
	var video models.Video
	pipeline := bson.A{
		bson.M{"$unwind": "$courses"},
		bson.M{"$unwind": "$courses.sections"},
		bson.M{"$unwind": "$courses.sections.videos"},
		bson.M{"$match": bson.M{"courses.sections.videos._id": videoId}},
		bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": []interface{}{"$$ROOT", bson.M{"video": bson.M{"_id": "$courses.sections.videos._id", "url": "$courses.sections.videos.url"}}}}}},
		bson.M{"$project": bson.M{"courses": 0}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {

		log.Printf("Error while aggregate videos %v", err.Error())
		return video, errors.New(shared.UNABLE_TO_GET_VIDEO)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		err := cursor.Decode(&videoResponse)
		if err != nil {
			return models.Video{}, errors.New(shared.UNABLE_TO_GET_VIDEO)
		}

	}
	video = videoResponse.Video
	return video, nil

}
func GetDetailVideo(ctx context.Context, collection *mongo.Collection, videoId primitive.ObjectID) (models.Video, error) {
	var videoResponse VideoQueryResponse
	var video models.Video
	pipeline := bson.A{
		bson.M{"$unwind": "$courses"},
		bson.M{"$unwind": "$courses.sections"},
		bson.M{"$unwind": "$courses.sections.videos"},
		bson.M{"$match": bson.M{"courses.sections.videos._id": videoId}},
		bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": []interface{}{"$$ROOT", bson.M{"video": bson.M{"_id": "$courses.sections.videos._id", "name": "$courses.sections.videos.name", "url": "$courses.sections.videos.url", "groups": "$courses.sections.videos.groups", "teacher_id": "$courses.sections.videos.teacher_id", "created_at": "$courses.sections.videos.created_at", "updated_at": "$courses.sections.videos.updated_at"}}}}}},
		bson.M{"$project": bson.M{"courses": 0}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {

		log.Printf("Error while aggregate videos %v", err.Error())
		return video, errors.New(shared.UNABLE_TO_GET_VIDEO)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		err := cursor.Decode(&videoResponse)
		if err != nil {
			return models.Video{}, errors.New(shared.UNABLE_TO_GET_VIDEO)
		}

	}
	video = videoResponse.Video
	return video, nil

}

func GetVideoFile(videoUrl string) (os.File, error) {
	dir, errD := utils.GetStorageFile("videos")
	if errD != nil {
		log.Printf("Error getting storage file: %v", errD.Error())
		return os.File{}, errors.New(shared.UNABLE_TO_GET_VIDEO)
	}
	videoPath := path.Join(dir, videoUrl)
	file, err := os.Open(videoPath)
	if err != nil {
		log.Printf("Error opening file: %v", err.Error())
		return os.File{}, errors.New(shared.UNABLE_TO_GET_VIDEO)
	}
	return *file, nil

}

func EditVideo(ctx context.Context, collection *mongo.Collection, video models.Video) error {
	update := bson.M{
		"$set": bson.A{
			bson.M{"courses.$[course].sections.$[section].videos.$[video].name": video.Name},
			bson.M{"courses.$[course].sections.$[section].videos.$[video].groups": video.Groups},
		},
	}
	arrayFilters := bson.A{
		bson.M{"course.sections.videos._id": video.ID},
		bson.M{"section.videos._id": video.ID},
		bson.M{"video._id": video.ID},
	}
	opts := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
	up := collection.FindOneAndUpdate(ctx, bson.M{"courses.sections.videos._id": video.ID}, update, opts)
	err := up.Err()
	if err != nil {
		log.Println("Error updating video: ", err)
		return errors.New(shared.FILE_NOT_UPDATED)
	}

	return nil
}

func DeleteVideo(ctx context.Context, collection *mongo.Collection, videoId primitive.ObjectID) error {
	filter := bson.D{{"courses.sections.videos._id", videoId}}
	update := bson.M{
		"$pull": bson.M{
			"courses.$[course].sections.$[section].videos": bson.M{"_id": videoId},
		},
	}
	arrayFilters := bson.A{
		bson.M{"course.sections.videos._id": videoId},
		bson.M{"section.videos._id": videoId},
	}
	opts := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
	rs := collection.FindOneAndUpdate(ctx, filter, update, opts)
	err := rs.Err()
	if err != nil {
		log.Printf("Error deleting video object: %v", err)
		return err

	}
	return nil
}

func DeletePhysicalVideo(videoUrl string) error {
	dir, err := utils.GetStorageFile("videos")
	if err != nil {
		log.Printf("Error getting storage file: %v", err.Error())
		return errors.New(shared.FILE_NOT_DELETED)
	}
	if videoUrl == "" {
		return errors.New(shared.FILE_NOT_DELETED)
	}
	videoPath := path.Join(dir, videoUrl)
	err = os.Remove(videoPath)
	if err != nil {
		log.Printf("Error deleting file object: %v", err.Error())
		return errors.New(shared.FILE_NOT_DELETED)
	}
	return nil
}
func OnCompleteVideo(ctx context.Context, collection *mongo.Collection, videoId primitive.ObjectID) (int32, error) {
	var video models.Video
	filter := bson.D{{"courses.sections.videos._id", videoId}}
	err := collection.FindOne(ctx, filter).Decode(&video)
	if err != nil {
		log.Printf("Error getting video object: %v", err.Error())
		return -1, errors.New(shared.UNABLE_TO_GET_VIDEO)
	}
	return video.Score, nil
}
