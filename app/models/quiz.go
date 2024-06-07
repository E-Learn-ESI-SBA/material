package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	ModuleId      string 			 `json:"module_id" bson:"module_id" validate:"required" binding:"required"`
	TeacherId     string             `json:"teacher_id" bson:"teacher_id"`
	Title         string             `json:"title" bson:"title" validate:"required" binding:"required"`
	Instructions  string             `json:"instructions" bson:"instructions" validate:"required" binding:"required"`
	Image 		  string			 `json:"image" bson:"image"`
	Year 		  string			 `json:"year" bson:"year" validate:"required" binding:"required"`
	QuestionCount int                `json:"question_count" bson:"question_count" validate:"required" binding:"required"`
	MaxScore      float64            `json:"max_score" bson:"max_score" validate:"required" binding:"required"`
	MinScore      float64            `json:"min_score" bson:"min_score" validate:"required" binding:"required"`
	StartDate     time.Time          `json:"start_date" bson:"start_date" validate:"required" binding:"required"`
	EndDate       time.Time          `json:"end_date" bson:"end_date" validate:"required" binding:"required"`
	Duration      int                `json:"duration" bson:"duration" validate:"required" binding:"required"`
	Questions     []Question         `json:"questions" bson:"questions" validate:"required" binding:"required"`
	Grades        []Grade            `json:"grades" bson:"grades" validate:"required" binding:"required"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
}

type QuizUpdate struct {
	Title 	   string             `json:"title" bson:"title" validate:"required" binding:"required"`
	Instructions string           `json:"instructions" bson:"instructions" validate:"required" binding:"required"`
	Image 	   string             `json:"image" bson:"image"`
	StartDate  time.Time          `json:"start_date" bson:"start_date" validate:"required" binding:"required"`
	EndDate    time.Time          `json:"end_date" bson:"end_date" validate:"required" binding:"required"`
	Duration   int                `json:"duration" bson:"duration" validate:"required" binding:"required"`
}

type Question struct {
	ID          string 				`json:"id" bson:"_id" validate:"required" binding:"required"`
	Body        string             `json:"body" bson:"body" validate:"required" binding:"required"`
	Score       float64            `json:"score" bson:"score" validate:"required" binding:"required"`
	Image       string             `json:"image" bson:"image"`
	Options     []Option           `json:"options" bson:"options" validate:"required" binding:"required"`
	CorrectIdxs []string              `json:"correct_idxs" bson:"correct_idxs" validate:"required" binding:"required"`
}

type Option struct {
	ID 		string 	`json:"id" bson:"id"`
	Option 	string 	`json:"option" bson:"option" validate:"required" binding:"required"` 
}

type Grade struct {
	Min   uint   `json:"min" bson:"min" validate:"required" binding:"required"`
	Max   uint   `json:"max" bson:"max" validate:"required" binding:"required"`
	Grade string `json:"grade" bson:"grade" validate:"required" binding:"required"`
}

type Answer struct {
	QuestionId string 				 `json:"question_id" bson:"question_id" validate:"required" binding:"required"`
	Choices    []string             `json:"choices" bson:"choices" validate:"required" binding:"required"`
	IsCorrect  bool              `json:"is_correct" bson:"is_correct"`
}
// choices will be compared against correctIdxs in the question
type Submission struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	StudentId string             `json:"student_id" bson:"student_id"`
	QuizId    primitive.ObjectID `json:"quiz_id" bson:"quiz_id"`
	Answers   []Answer           `json:"answers" bson:"answers" validate:"required" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	Grade     string             `json:"grade" bson:"grade"`
	Score     float64            `json:"score" bson:"score"`
	IsPassed  bool               `json:"is_passed" bson:"is_passed"`
}
