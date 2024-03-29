package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID string `json:"id"`
}
type Module struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" validate:"" `
	Year        int8               `json:"year"`
	Semester    int8               `json:"semester"`
	Coefficient int8               `json:"coefficient"`
	Courses     []Course           `json:"courses" bson:"courses"`
	TeacherId   string             `json:"teacher_id" bson:"teacher_id" validate:"required"`
}

type Course struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Plan        []string           `json:"plan" validate:"required,min=1" bson:"plan"`
	Sections    []Section          `json:"sections" bson:"sections"`
}

type Section struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" validate:"required" `
	CourseID string             `json:"course_id" validate:"required"`
	Lectures []Lecture          `json:"lectures"`
	Videos   []Video            `json:"videos"`
}

type Lecture struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Name    string             `json:"name"`
	Content string             `json:"content" bson:"content" validate:"required,min=250"`
}
type Video struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Url         string             `json:"url" bson:"url" validate:"required"`
	IsCompleted bool               `json:"isCompleted" bson:"isCompleted"`
}

type Years struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Year       int8               `json:"year"`
	Speciality string             `json:"speciality"`
}
