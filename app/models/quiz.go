package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Quiz struct {
	ID       		primitive.ObjectID 	`json:"id" bson:"_id"`
	ModuleId 		primitive.ObjectID 	`json:"module_id" bson:"module_id" validate:"required" binding:"required"`
	TeacherId 		int 				`json:"teacher_id" bson:"teacher_id"`
	Title    		string             	`json:"title" bson:"title" validate:"required" binding:"required"`
	Instructions 	string         		`json:"instructions" bson:"instructions" validate:"required" binding:"required"`
	QuestionCount 	int          		`json:"question_count" bson:"question_count" validate:"required" binding:"required"`
	MaxScore 		float64         	`json:"max_score" bson:"max_score" validate:"required" binding:"required"`
	StartDate 		time.Time 	   		`json:"start_date" bson:"start_date" validate:"required" binding:"required"`
	EndDate   		time.Time         	`json:"end_date" bson:"end_date" validate:"required" binding:"required"`
	Duration 		int               	`json:"duration" bson:"duration" validate:"required" binding:"required"`
	Questions 		[]Question 			`json:"questions" bson:"questions" validate:"required" binding:"required"`
	Grades 			[]Grade 			`json:"grades" bson:"grades" validate:"required" binding:"required"`
	Date
}

type Question struct {
	ID 			primitive.ObjectID 	`json:"id" bson:"_id"`
	Body 		string 				`json:"body" bson:"body" validate:"required" binding:"required"`
	Description string 				`json:"description" bson:"description" validate:"required"`
	Score 		float64 			`json:"score" bson:"score" validate:"required" binding:"required"`
	Image 		string 				`json:"image" bson:"image"`
	Options 	[]string 			`json:"options" bson:"options" validate:"required" binding:"required"`
	CorrectIdxs []int				`json:"correct_idxs" bson:"correct_idxs" validate:"required" binding:"required"`
}


type Grade struct {
	Min 		uint `json:"min" bson:"min" validate:"required" binding:"required"`
	Max 		uint `json:"max" bson:"max" validate:"required" binding:"required"`
	Grade 		string `json:"grade" bson:"grade" validate:"required" binding:"required"`
	IsPass 		bool `json:"is_pass" bson:"is_pass" validate:"required" binding:"required"`
}

type Answer struct {
	QuestionId 	primitive.ObjectID 	`json:"question_id" bson:"question_id" validate:"required" binding:"required"`
	Choices 	[]int 			`json:"choices" bson:"choices" validate:"required" binding:"required"`
	IsCorrect 	bool 				`json:"is_correct" bson:"is_correct"`
}

type Submission struct {
	ID 			primitive.ObjectID 	`json:"id" bson:"_id"`
	StudentId 	int 				`json:"student_id" bson:"student_id"`
	QuizId 		primitive.ObjectID 	`json:"quiz_id" bson:"quiz_id"`
	Answers 	[]Answer 			`json:"answers" bson:"answers" validate:"required" binding:"required"`
	CreatedAt 	time.Time 			`json:"created_at" bson:"created_at"`
	Grade 		string 				`json:"grade" bson:"grade"`
	Score 		float64 			`json:"score" bson:"score"`
}