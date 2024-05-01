package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
)

func GetCourseCommentsByCourseId(ctx context.Context, collection *mongo.Collection, courseId string, pagination interfaces.PaginationQuery) ([]models.Comments, error) {
	var comments []models.Comments
	filter := bson.D{{"course_id", courseId}}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetSkip(int64(pagination.Page * pagination.Items)).SetLimit(int64(pagination.Items))
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Error While Getting the Comments: %v\n", err)
		return nil, err

	}
	cursorError := cursor.All(ctx, &comments)
	if cursorError != nil {
		log.Printf("Error While Parsing Comments: %v\n", cursorError)
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

/*
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
*/
func EditComment(ctx context.Context, collection *mongo.Collection, commentId string, userId string, editedComment models.Comments) error {
	var comment models.Comments
	id, errId := primitive.ObjectIDFromHex(commentId)
	if errId != nil {
		log.Printf("Error While Parsing Section ID: %v\n", errId)
		return errId
	}
	filter := bson.D{{"_id", id}, {"user_id", userId}}
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
func EditReply(ctx context.Context, collection *mongo.Collection, commentId string, replyId string, userId string, editedReply models.Reply) error {
	var reply models.Reply
	id, errId := primitive.ObjectIDFromHex(commentId)
	if errId != nil {
		log.Printf("Error While Parsing Section ID: %v\n", errId)
		return errId
	}
	id2, errId2 := primitive.ObjectIDFromHex(replyId)
	if errId2 != nil {
		log.Printf("Error While Parsing Section ID: %v\n", errId)
		return errId2
	}
	// find replay from the comment.replays array
	filter := bson.D{{"_id", id}, {"replays._id", id2}, {"res.user_id", userId}}
	update := bson.D{{
		"$set", bson.D{{"replays.$.content", editedReply.Content}},
	}}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&reply)
	if err != nil {
		log.Printf("Error While Updating the Reply: %v\n", err)
		return err
	}
	return nil
}
func DeleteCommentByUser(ctx context.Context, collection *mongo.Collection, commentId string, userId string) error {
	id, errId := primitive.ObjectIDFromHex(commentId)
	if errId != nil {
		log.Printf("Error While Parsing Section ID: %v\n", errId)
		return errId
	}
	filter := bson.D{{"_id", id}, {"user_id", userId}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error While Deleting the Comment: %v\n", err)
		return err
	}
	return nil
}
func DeleteReplyByUser(ctx context.Context, collection *mongo.Collection, commentId string, replyId string, userId string) error {
	id, errId := primitive.ObjectIDFromHex(commentId)
	if errId != nil {
		log.Printf("Error While Parsing Section ID: %v\n", errId)
		return errId
	}
	id2, errId2 := primitive.ObjectIDFromHex(replyId)
	if errId2 != nil {
		log.Printf("Error While Parsing Section ID: %v\n", errId)
		return errId2
	}
	// Remove replay from the comment.replays array
	filter := bson.D{{"_id", id}, {"replays._id", id2}, {"user_id", userId}}
	update := bson.D{{
		"$pull", bson.D{{"replays", bson.D{{"_id", replyId}}}},
	}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error While Deleting the Replay: %v\n", err)
	}
	return err

}
func ReplayToComment(ctx context.Context, collection *mongo.Collection, replay models.Reply, commentId primitive.ObjectID) error {
	var comment models.Comments
	filter := bson.D{{"_id", commentId}}
	// insert into replays array
	// before insert make sure the replays under 10 replays
	/* const maxReplays = 10
	documents, err := collection.CountDocuments(ctx, bson.D{{"_id", commendId}, {"replays", bson.D{{"$size", maxReplays}}}})
	if err != nil {
		return err
	}
	if documents >= maxReplays {

		return errors.New("Replays are Full,  Can't add more than 10 replays")

	}
	*/
	update := bson.D{{
		"$push", bson.D{{"replays", replay}}},
	}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&comment)
	if err != nil {
		log.Printf("Error While Updating the Comment: %v\n", err)
		return err
	}
	return nil

}
func CreateComment(ctx context.Context, collection *mongo.Collection, comment models.Comments) error {
	_, err := collection.InsertOne(ctx, comment)
	if err != nil {
		log.Printf("Error While Creating the Comment: %v\n", err)
		return err
	}
	return nil
}
