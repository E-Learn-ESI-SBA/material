package services

import (
	"context"
	"errors"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
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


	quiz.ID = primitive.NewObjectID()
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
	quizUpdates models.QuizUpdate,
	quizID string,
	teacherID string,
) error {
	objectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return err
	}

	log.Printf("quizUpdates: %v\n", quizUpdates)
	log.Printf("here")
	filter := bson.D{{"_id", objectID}, {"teacher_id", teacherID}}
	// should be updated field by field to avoid overriding existing data with nulls
	updates := bson.D{{"$set", bson.D{
		{"title", quizUpdates.Title},
		{"instructions", quizUpdates.Instructions},
		{"start_date", quizUpdates.StartDate},
		{"end_date", quizUpdates.EndDate},
		{"duration", quizUpdates.Duration},
		{"image", quizUpdates.Image},
		{"updated_at", time.Now()},
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
	submissionsCollection *mongo.Collection,
	quizID string,
	studentId string,
) (models.Quiz, error, bool) {
	var quiz models.Quiz
	objectId, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return quiz, err, false
	}
	filter := bson.D{{"_id", objectId}}
	projection := bson.D{
		{"questions", 0},
		{"grades", 0},
	}
	err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&quiz)
	if err != nil {
		log.Printf("Error While Getting Quiz: %v\n", err)
		return quiz, err, false
	}

	// check if submissuion exists
	var submission models.Submission
	filter = bson.D{{"quiz_id", objectId}, {"student_id", studentId}}
	err = submissionsCollection.FindOne(ctx, filter).Decode(&submission)
	log.Printf("submission: %v\n", submission)
	log.Printf(("quiz id: %v\n"), quiz.ID)
	log.Printf("student id: %v\n", studentId)
	if (submission.ID.IsZero()) {
		return quiz, nil, false
	}

	return quiz, nil, true
}

func GetManyQuizesByTeacherId(
	ctx context.Context,
	collection *mongo.Collection,
	teacherID string,
) ([]models.Quiz, error) {
	var quizes []models.Quiz
	cursor, err := collection.Find(ctx, bson.M{"teacher_id": teacherID})
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


func GetQuizesByStudentId(
	ctx context.Context,
	collection *mongo.Collection,
	student *utils.UserDetails,
) ([]models.Quiz, error) {
	var quizes []models.Quiz

	pipeline := mongo.Pipeline{
		{{
			Key: "$match",
			Value: bson.D{
				{Key: "start_date", Value: bson.D{{Key: "$lte", Value: time.Now()}}},
				{Key: "year", Value: student.Year},
			},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var quiz models.Quiz
		if err := cursor.Decode(&quiz); err != nil {
			return nil, err
		}
		quizes = append(quizes, quiz)
	}

	// Check if there was an error during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return quizes, nil
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
) (int32, error) {

	// check if the quiz exists
	var quiz models.Quiz
	filter := bson.D{{"_id", submission.QuizId}}
	err := collection.FindOne(ctx, filter).Decode(&quiz)
	if err != nil {
		log.Printf("Error While Getting Quiz: %v\n", err)
		return 0, err
	}

	if time.Now().Before(quiz.StartDate) || time.Now().After(quiz.EndDate) {
		return 0, errors.New("quiz is not ongoing")
	}
	log.Printf("quiz id %v\n", quiz.ID)
	log.Printf("student id %v\n", submission.StudentId)

	// check if the student already submitted the quiz answers
	filter = bson.D{{"quiz_id", submission.QuizId}, {"student_id", submission.StudentId}}
	var existingSubmission models.Submission
	err = SubmissionsCollection.FindOne(ctx, filter).Decode(&existingSubmission)
	if err == nil {
		return 0, errors.New(shared.QUIZ_ANSWER_ALREADY_SUBMITTED)
	}

	questions := quiz.Questions
	submission.ID = primitive.NewObjectID()
	submission.CreatedAt = time.Now()
	submission.Score, submission.Answers = CalcFinalScore(questions, submission.Answers)
	submission.Grade = GetGrade(submission.Score, quiz.Grades)
	if submission.Score >= quiz.MinScore {
		submission.IsPassed = true
	}
	_, err = SubmissionsCollection.InsertOne(ctx, submission)
	if err != nil {
		log.Printf("Error While Submitting Quiz Answers: %v\n", err)
		return 0, err
	}

	return int32(submission.Score), nil
}

func GetQuizResultByStudentId(
	ctx context.Context,
	collection *mongo.Collection,
	submissionsCollection *mongo.Collection,
	moduleCollection *mongo.Collection,
	quizID string,
	studentID string,
) (models.Submission, models.Quiz, *string, error) {
	objectId, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return models.Submission{}, models.Quiz{}, nil, err
	}
	// check if the quiz exists
	var quiz models.Quiz
	filter := bson.D{{"_id", objectId}}
	_ = collection.FindOne(ctx, filter).Decode(&quiz)
	if quiz.ID.IsZero() {
		log.Printf("Quiz Not Found\n")
		return models.Submission{}, models.Quiz{}, nil, errors.New(shared.QUIZ_NOT_FOUND)
	}
	// check if quiz has ended
	if time.Now().Before(quiz.EndDate) {
		return models.Submission{}, models.Quiz{}, nil, errors.New(shared.QUIZ_STILL_ONGOING)
	}
	filter = bson.D{{"quiz_id", objectId}, {"student_id", studentID}}
	var submission models.Submission
	err = submissionsCollection.FindOne(ctx, filter).Decode(&submission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return submission, models.Quiz{}, nil, errors.New(shared.QUIZ_ANSWER_NOT_FOUND)
		}
		return submission, models.Quiz{}, nil, err
	}

	// get module name
	var module models.Module
	options := options.FindOne().SetProjection(bson.M{"name": 1})
	_ = moduleCollection.FindOne(ctx, bson.M{"_id": quiz.ModuleId}, options).Decode(&module)
	var module_name string
	if module.ID.IsZero() {
		module_name = "placeholder"
		return submission, quiz, &module_name, nil
	}

	return submission, quiz, &module.Name, nil
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
		bson.D{{"$project", bson.D{
			{"answers", 0}, // Exclude the answers field
		}}},
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
	submissionCollection *mongo.Collection,
	quizID string,
	studentID string,
) (models.Quiz, error) {
	objectId, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return models.Quiz{}, err
	}
	// check if submissuion exists
	var submission models.Submission
	filter := bson.D{{"quiz_id", objectId}, {"student_id", studentID}}
	_ = submissionCollection.FindOne(ctx, filter).Decode(&submission)
	if (!submission.ID.IsZero()) {
		return models.Quiz{}, errors.New(shared.QUIZ_ANSWER_ALREADY_SUBMITTED)
	}
	// check if the quiz exists
	var quiz models.Quiz
	_ = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&quiz)
	if quiz.ID.IsZero() {
		return models.Quiz{}, errors.New("quiz not found")
	}


	//check if start date is before now
	if time.Now().Before(quiz.StartDate) {
		return models.Quiz{}, errors.New(shared.QUIZ_NOT_STARTED)
	}

	
	return quiz, nil
}


func GetSubmissionDetails(
	ctx context.Context,
	collection *mongo.Collection,
	moduleCollection *mongo.Collection,
	submissionsCollection *mongo.Collection,
	submissionID string,
) (models.Submission, *string, models.Quiz, error) {
	objectId, err := primitive.ObjectIDFromHex(submissionID)
	if err != nil {
		log.Printf("Error While Converting ID: %v\n", err)
		return models.Submission{}, nil, models.Quiz{}, err
	}
	
	var submission models.Submission
	_ = submissionsCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&submission)
	if submission.ID.IsZero() {
		return models.Submission{}, nil, models.Quiz{}, errors.New("submission not found")
	}

	// fetch quiz
	var quiz models.Quiz
	_ = collection.FindOne(ctx, bson.M{"_id": submission.QuizId}).Decode(&quiz)
	if quiz.ID.IsZero() {
		return models.Submission{}, nil, models.Quiz{}, errors.New("quiz not found")
	}

	// fetch module name return placeholder if module not found
	var module models.Module
	options := options.FindOne().SetProjection(bson.M{"name": 1})
	_ = collection.FindOne(ctx, bson.M{"_id": quiz.ModuleId}, options).Decode(&module)
	var module_name string
	if module.ID.IsZero() {
		module_name = "placeholder"
		return submission, &module_name, quiz, nil
	}

	return submission, &module.Name, quiz, nil
}


func CalcFinalScore(questions []models.Question, answers []models.Answer) (float64, []models.Answer) {
	var totalScore float64
	// i have 0 brain cells left
	for _, question := range questions {
	 for i, _ := range answers {
	  if question.ID == answers[i].QuestionId {
	   if AllIn(question.CorrectIdxs, answers[i].Choices) && len(question.CorrectIdxs) == len(answers[i].Choices) {
		totalScore += question.Score
		answers[i].IsCorrect = true
	   } else {
		answers[i].IsCorrect = false
	   }
	   break
	  }
	 }
	}
	return totalScore, answers
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
