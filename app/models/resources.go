package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Section struct {
	ID
	Name      string     `json:"name" validate:"required" `
	CourseID  string     `json:"course_id" validate:"required"`
	Order     int8       `json:"order" bson:"order" validate:"min=1"`
	TeacherId int        `json:"teacher_id" bson:"teacher_id" validate:"required"`
	Videos    *[]Video   `json:"videos" bson:"videos"`
	Lectures  *[]Lecture `json:"lectures" bson:"lectures"`
	Files     *[]Files   `json:"files" bson:"files"`
	Date
}

type Lecture struct {
	ID
	Name      string             `json:"name"`
	Content   string             `json:"content" bson:"content" validate:"required,min=250"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required"`
	TeacherId int                `json:"teacher_id" bson:"teacher_id" validate:"required"`
	IsPublic  bool               `json:"is_public" bson:"is_public" validate:"default=false"`
	Date
}
type Video struct {
	ID
	Url       string             `json:"url" bson:"url" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required"`
	TeacherId int                `json:"teacher_id" bson:"teacher_id" validate:"required"`
	Date
}
type Files struct {
	ID
	Url       string             `json:"url" bson:"url" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required"`
	TeacherId int                `json:"teacher_id" bson:"teacher_id" validate:"required"`
	Group     int8               `json:"group" bson:"group" validate:"required"`
	Date
}

type StudentNote struct {
	ID
	StudentID int                `json:"student_id" bson:"student_id" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required"`
	Content   string             `json:"content" bson:"content" validate:"required"`
}
