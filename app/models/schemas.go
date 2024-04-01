package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"madaurus/dev/material/app/utils"
	"time"
)

type User struct {
	ID string `json:"id"`
}

type Date struct {
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
type Module struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" validate:"" `
	Year        int8               `json:"year"`
	Speciality  *string            `json:"speciality,omitempty"`
	Semester    int8               `json:"semester"`
	Coefficient int8               `json:"coefficient"`
	TeacherId   int                `json:"teacher_id" bson:"teacher_id" validate:"required"`
	IsPublic    bool               `json:"isPublic" bson:"isPublic" validate:"default=false"`
	Image       *string            `json:"image,omitempty"`
	User        utils.LightUser    `json:"user"`
	Date
}

type Course struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Plan        []string           `json:"plan" validate:"required,min=1" bson:"plan"`
	ModuleId    primitive.ObjectID `json:"module_id" bson:"module_id"`
	Sections    *[]Section         `json:"sections" bson:"sections"`
	Date
}

type Section struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" validate:"required" `
	CourseID string             `json:"course_id" validate:"required"`
	Order    int8               `json:"order" bson:"order" validate:"min=1"`
	Videos   *[]Video           `json:"videos" bson:"videos"`
	Lectures *[]Lecture         `json:"lectures" bson:"lectures"`
	Files    *[]Files           `json:"files" bson:"files"`
	Date
}

type Lecture struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name"`
	Content   string             `json:"content" bson:"content" validate:"required,min=250"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required"`
	Date
}
type Video struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Url       string             `json:"url" bson:"url" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required"`
	Date
}
type Files struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Url       string             `json:"url" bson:"url" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required"`
	TeacherId int                `json:"teacher_id" bson:"teacher_id" validate:"required"`
	Group     int8               `json:"group" bson:"group" validate:"required"`
	Date
}

type Comments struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Content  string             `json:"content" bson:"content" validate:"required"`
	UserId   int                `json:"user_id" bson:"user_id" validate:"required"`
	CourseId primitive.ObjectID `json:"course_id" bson:"course_id" validate:"required"`
	Date
	Replays *[]Reply `json:"replays" bson:"replays"`
}

type Reply struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Content string             `json:"content" bson:"content" validate:"required"`
	UserId  string             `json:"user_id" bson:"user_id" validate:"required"`
	User    utils.LightUser    `json:"user" bson:"user"`
	Date
}

type StudentNote struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	StudentID int                `json:"student_id" bson:"student_id" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required"`
	Content   string             `json:"content" bson:"content" validate:"required"`
}
