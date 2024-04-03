package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"madaurus/dev/material/app/models"
)

func DeleteStudentNote(collection *mongo.Collection, ctx context.Context, noteId string) error {
	filter := bson.D{{"_id", noteId}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error While Deleting Student Note: %v\n", err)
		return err
	}
	return nil
}

func CreateStudentNote(collection *mongo.Collection, ctx context.Context, note models.StudentNote) error {
	_, err := collection.InsertOne(ctx, note)
	if err != nil {
		log.Printf("Error While Creating Student Note: %v\n", err)
		return err
	}
	return nil
}

func GetStudentNotes(ctx context.Context, collection *mongo.Collection, studentId string, sectionId string) (models.StudentNote, error) {
	var note models.StudentNote
	filter := bson.D{{"student_id", studentId}, {"section_id", sectionId}}
	// select only the title
	opts := options.Find().SetProjection(bson.D{{"title", 1}, {"_id", 1}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Error While Getting Student Note: %v\n", err)
		return note, err
	}
	cursorError := cursor.All(ctx, &note)
	if cursorError != nil {
		log.Printf("Error While Parsing Student Note: %v\n", cursorError)
		return note, cursorError

	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("failed to close cursor")
		}
	}(cursor, ctx)
	return note, nil
}

func GetNoteById(ctx context.Context, collection *mongo.Collection, noteId string) (models.StudentNote, error) {
	var note models.StudentNote
	filter := bson.D{{"_id", noteId}}
	err := collection.FindOne(ctx, filter).Decode(&note)
	if err != nil {
		log.Printf("Error While Getting Note: %v\n", err)
		return note, err
	}
	return note, nil
}

func UpdateStudentNote(ctx context.Context, collection *mongo.Collection, note models.StudentNote) error {
	filter := bson.D{{"_id", note.ID}}
	update := bson.D{{"$set", note}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error While Updating Student Note: %v\n", err)
		return err
	}
	return nil
}
