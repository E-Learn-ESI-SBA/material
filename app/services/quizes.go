package services

import (
	"context"
	"errors"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateQuiz(
	ctx context.Context,
	collection *mongo.Collection,
	moduleCollection *mongo.Collection,
	quiz models.Quiz,
) error {
	var module models.Module
	moduleObjectID, err := primitive.ObjectIDFromHex(quiz.ModuleId)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return err
	}
	filter := bson.D{{"_id", moduleObjectID}, {"teacher_id", quiz.TeacherId}}
	moduleCollection.FindOne(ctx, filter).Decode(&module)

	// will uncomment this when i fix permit io
	// if module.ID.IsZero() {
	// 	return errors.New("module not found")
	// }


	quiz.CreatedAt = time.Now()
	quiz.UpdatedAt = quiz.CreatedAt
	_, err = collection.InsertOne(ctx, quiz)
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
	teacherID string,
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
		return errors.New(shared.UNAUTHORIZED)
	}
	return nil
}

func DeleteQuiz(
	ctx context.Context,
	collection *mongo.Collection,
	quizID string,
	teacherID string,
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
		return errors.New(shared.UNAUTHORIZED)
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

func GetQuizResults(
	ctx context.Context,
	collection *mongo.Collection,
	submissionsCollection *mongo.Collection,
	quizID string,
	teacherId string,
) ([]models.Submission, error) {
	objectId, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return nil, err
	}
	// check if the quiz exists
	var quiz models.Quiz
	filter := bson.D{{"_id", objectId}, {"teacher_id", teacherId}}
	_ = collection.FindOne(ctx, filter).Decode(&quiz)
	if quiz.ID.IsZero() {
		return nil, errors.New(shared.QUIZ_NOT_FOUND)
	}
	filter = bson.D{{"quiz_id", objectId}}
	var submissions []models.Submission
	cursor, err := submissionsCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error While Getting Submissions: %v\n", err)
		return submissions, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var submission models.Submission
		cursor.Decode(&submission)
		submissions = append(submissions, submission)
	}

	return submissions, nil
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
		return errors.New(shared.QUIZ_ANSWER_ALREADY_SUBMITTED)
	}

	questions := quiz.Questions
	submission.CreatedAt = time.Now()
	submission.Score = CalcFinalScore(questions, submission.Answers)
	submission.Grade = GetGrade(submission.Score, quiz.Grades)
	if submission.Score >= quiz.MinScore {
		submission.IsPassed = true
	}
	_, err = SubmissionsCollection.InsertOne(ctx, submission)
	if err != nil {
		log.Printf("Error While Submitting Quiz Answers: %v\n", err)
		return err
	}

	return nil
}

func GetQuizResultByStudentId(
	ctx context.Context,
	collection *mongo.Collection,
	submissionsCollection *mongo.Collection,
	quizID string,
	studentID string,
) (models.Submission, error) {
	objectId, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return models.Submission{}, err
	}
	// check if the quiz exists
	var quiz models.Quiz
	filter := bson.D{{"_id", objectId}}
	_ = collection.FindOne(ctx, filter).Decode(&quiz)
	if quiz.ID.IsZero() {
		return models.Submission{}, errors.New(shared.QUIZ_NOT_FOUND)
	}
	// check if quiz has ended
	if time.Now().Before(quiz.EndDate) {
		return models.Submission{}, errors.New(shared.QUIZ_STILL_ONGOING)
	}
	filter = bson.D{{"quiz_id", objectId}, {"student_id", studentID}}
	var submission models.Submission
	err = submissionsCollection.FindOne(ctx, filter).Decode(&submission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return submission, errors.New(shared.QUIZ_ANSWER_NOT_FOUND)
		}
		return submission, err
	}
	return submission, nil
}

func GetQuizesResultsByStudentId(
	ctx context.Context,
	collection *mongo.Collection,
	submissionsCollection *mongo.Collection,
	studentID string,
) ([]models.Submission, error) {
	var submissions []models.Submission
	// fetch submissions by student id where time.Now() > quiz.end_date
	pipe := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"student_id", studentID}}}},
		bson.D{{"$lookup", bson.D{
			{"from", "quizes"},
			{"localField", "quiz_id"},
			{"foreignField", "_id"},
			{"as", "quiz"},
		}}},
		bson.D{{"$unwind", "$quiz"}},
		bson.D{{"$match", bson.D{{"quiz.end_date", bson.D{{"$lt", time.Now()}}}}}},
	}
	cursor, err := submissionsCollection.Aggregate(ctx, pipe)
	if err != nil {
		log.Printf("Error While Getting Submissions: %v\n", err)
		return submissions, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var submission models.Submission
		cursor.Decode(&submission)
		submissions = append(submissions, submission)
	}
	
	return submissions, nil
}

func GetQuizQuestions(
	ctx context.Context,
	collection *mongo.Collection,
	quizID string,
	studentID string,
) ([]models.Question, error) {
	objectId, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return []models.Question{}, err
	}
	// check if the quiz exists
	var quiz models.Quiz
	opts := options.FindOne().SetProjection(bson.M{"questions": 1})
	_ = collection.FindOne(ctx, bson.M{"_id": objectId}, opts).Decode(&quiz)
	if quiz.ID.IsZero() {
		return nil, errors.New("quiz not found")
	}

	//check if start date is before now
	if time.Now().Before(quiz.StartDate) {
		return nil, errors.New(shared.QUIZ_NOT_STARTED)
	}

	questions := quiz.Questions
	
	return questions, nil
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

func isIn(val string, arr []string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func AllIn(arr1 []string, arr2 []string) bool {
	for _, v := range arr1 {
		if !isIn(v, arr2) {
			return false
		}
	}
	return true
}
