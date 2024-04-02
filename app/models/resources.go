package models

import "go.mongodb.org/mongo-driver/bson/primitive"

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

type StudentNote struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	StudentID int                `json:"student_id" bson:"student_id" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required"`
	Content   string             `json:"content" bson:"content" validate:"required"`
}
