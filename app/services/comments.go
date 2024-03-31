package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/models"
)

func GetCourseCommentsByCourseId(ctx context.Context, collection *mongo.Collection, courseId string) ([]models.Comments, error) {
	var comments []models.Comments
	cursor, err := collection.Find(ctx, bson.D{{"course_id", courseId}})
	if err != nil {
		return nil, err
	}
	cursorError := cursor.All(ctx, &comments)
	if cursorError != nil {
		return nil, cursorError
	}
	return comments, nil
}

func GetCommentReplay(ctx context.Context, collection *mongo.Collection, commentId string) ([]models.Reply, error) {
	var comments []models.Reply
	filter := bson.D{{"comment_id", commentId}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error While Getting the Replays: %v\n", err)
		return nil, err
	}
	cursorError := cursor.All(ctx, &comments)
	if cursorError != nil {
		log.Printf("Error While Parsing   Replays: %v\n", cursorError)
		return nil, cursorError
	}

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}()
	return comments, nil
}

func EditComment(ctx context.Context, collection *mongo.Collection, commentId string, userId int, editedComment models.Comments) error {
	var comment models.Comments
	filter := bson.D{{"_id", commentId}, {"user_id", userId}}
	update := bson.D{{
		"$set", bson.D{{"content", editedComment.Content}},
	}}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&comment)
	if err != nil {
		log.Printf("Error While Updating the Comment: %v\n", err)
		return err
	}
	return nil

}
func EditReply(ctx context.Context, collection *mongo.Collection, replyId string, userId int, editedReply models.Reply) error {
	var reply models.Reply
	filter := bson.D{{"_id", replyId}, {"user_id", userId}}
	update := bson.D{{
		"$set", bson.D{{"content", editedReply.Content}},
	}}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&reply)
	if err != nil {
		log.Printf("Error While Updating the Reply: %v\n", err)
		return err
	}
	return nil
}
func DeleteCommentByUser(ctx context.Context, collection *mongo.Collection, commentId string, userId int) error {
	filter := bson.D{{"_id", commentId}, {"user_id", userId}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error While Deleting the Comment: %v\n", err)
		return err
	}
	return nil
}
func DeleteReplyByUser(ctx context.Context, collection *mongo.Collection, replyId string, userId int) error {
	filter := bson.D{{"_id", replyId}, {"user_id", userId}}
	var reply models.Reply
	err := collection.FindOneAndDelete(ctx, filter).Decode(&reply)
	if err != nil {
		log.Printf("Error While Deleting the Reply: %v\n", err)
		return err
	}
	return nil

}
