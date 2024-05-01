package services

import (
	"context"
	"errors"
	"log"
	"madaurus/dev/material/app/models"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



func CreateQuiz(
	ctx context.Context, 
	collection *mongo.Collection, 
	moduleCollection *mongo.Collection, 
	quiz models.Quiz,
	) error {
	var module models.Module
	filter := bson.D{{"_id", quiz.ModuleId}, {"teacher_id", quiz.TeacherId}}
	moduleCollection.FindOne(ctx, filter).Decode(&module)
	
	if module.ID.IsZero() {
		return errors.New("module not found")
	}
	
	currTime := time.Now()
	quiz.Date.CreatedAt = &currTime
	_, err := collection.InsertOne(ctx, quiz)
	if err != nil {
		log.Printf("Error While Creating Quiz: %v\n", err)
		return err
	}
	return nil
}


func UpdateQuiz(
	ctx context.Context, 
	collection *mongo.Collection, 
	quiz models.Quiz, 
	teacherID int,
	) error {
	filter := bson.D{{"_id", quiz.ID}, {"teacher_id", teacherID}}
	// should be updated field by field to avoid overriding existing data with nulls
	updates := bson.D{{"$set", bson.D{
		{"title", quiz.Title},
		{"instructions", quiz.Instructions},
		{"question_count", quiz.QuestionCount},
		{"start_date", quiz.StartDate},
		{"end_date", quiz.EndDate},
		{"duration", quiz.Duration},
		{"date.updated_at", time.Now()},
	}}}
	res, err := collection.UpdateOne(ctx, filter, updates)
	if err != nil {
		log.Printf("Error While Updating Quiz: %v\n", err)
		return err
	} else if res.MatchedCount == 0 {
		return errors.New("Unauthorized")
	}
	return nil
}


func DeleteQuiz(
	ctx context.Context, 
	collection *mongo.Collection, 
	quizID string, 
	teacherID int,
	) error {
	objectId, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return err
	}
	filter := bson.D{{"_id", objectId}, {"teacher_id", teacherID}}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error While Deleting Course: %v\n", err)
		return err
	} else if res.DeletedCount == 0 {
		return errors.New("Unauthorized")
	}
	return nil
}


func GetQuiz(
	ctx context.Context,
	collection *mongo.Collection,
	quizID string,
	) (models.Quiz, error) {
	var quiz models.Quiz
	objectId, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return quiz, err
	}
	filter := bson.D{{"_id", objectId}}
	err = collection.FindOne(ctx, filter).Decode(&quiz)
	if err != nil {
		log.Printf("Error While Getting Quiz: %v\n", err)
		return quiz, err
	}
	return quiz, nil
}


func GetQuizesByModuleId(
	ctx context.Context,
	collection *mongo.Collection,
	moduleID string,
	) ([]models.Quiz, error) {
	var quizes []models.Quiz
	objectId, err := primitive.ObjectIDFromHex(moduleID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return quizes, err
	}
	cursor, err := collection.Find(ctx, bson.M{"module_id": objectId})
	if err != nil {
		log.Printf("Error While Getting Quizes: %v\n", err)
		return quizes, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var quiz models.Quiz
		cursor.Decode(&quiz)
		quizes = append(quizes, quiz)
	}
	return quizes, nil
}


func GetQuizesByAdmin(
	ctx context.Context,
	collection *mongo.Collection,
	) ([]models.Quiz, error) {
	var quizes []models.Quiz
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error While Getting Quizes: %v\n", err)
		return quizes, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var quiz models.Quiz
		cursor.Decode(&quiz)
		quizes = append(quizes, quiz)
	}
	return quizes, nil
}


func SubmitQuizAnswers(
	ctx context.Context,
	collection *mongo.Collection,
	SubmissionsCollection *mongo.Collection,
	submission models.Submission,
	) error {
	
	// check if the quiz exists
	var quiz models.Quiz
	filter := bson.D{{"_id", submission.QuizId}}
	err := collection.FindOne(ctx, filter).Decode(&quiz)
	if err != nil {
		log.Printf("Error While Getting Quiz: %v\n", err)
		return err
	}

	if time.Now().Before(quiz.StartDate) || time.Now().After(quiz.EndDate) {
		return errors.New("quiz is not ongoing")
	}

	// check if the student already submitted the quiz answers
	filter = bson.D{{"quiz_id", submission.QuizId}, {"student_id", submission.StudentId}}
	var existingSubmission models.Submission
	err = SubmissionsCollection.FindOne(ctx, filter).Decode(&existingSubmission)
	if err == nil {
		return errors.New("student already submitted the quiz answers")
	}

	questions := quiz.Questions
	submission.CreatedAt = time.Now()
	submission.Score = CalcFinalScore(questions, submission.Answers)
	submission.Grade = GetGrade(submission.Score, quiz.Grades)
	_, err = SubmissionsCollection.InsertOne(ctx, submission)
	if err != nil {
		log.Printf("Error While Submitting Quiz Answers: %v\n", err)
		return err
	}

	return nil
}


func CalcFinalScore(questions []models.Question, answers []models.Answer) float64 {
	var totalScore float64
	// i have 0 brain cells left
	for _, question := range questions {
		for _, answer := range answers {
			if question.ID == answer.QuestionId {
				if AllIn(question.CorrectIdxs, answer.Choices) && len(question.CorrectIdxs) == len(answer.Choices) {
					totalScore += question.Score
					answer.IsCorrect = true
				} else {
					answer.IsCorrect = false
				}
				break
			}
		}
	}
	return totalScore
}

func GetGrade(score float64, grades []models.Grade) string {
	for _, grade := range grades {
		if uint(math.Round(score)) >= grade.Min && uint(math.Round(score)) <= grade.Max {
			return grade.Grade
		}
	}
	return ""
}


func isIn(val int, arr []int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func AllIn(arr1 []int, arr2 []int) bool {
	for _, v := range arr1 {
		if !isIn(v, arr2) {
			return false
		}
	}
	return true
}