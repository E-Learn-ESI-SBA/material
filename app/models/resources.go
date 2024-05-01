package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Section struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" validate:"required" binding:"required" bson:"name" `
	CourseID  string             `json:"course_id" validate:"required" binding:"required" bson:"course_id"`
	Order     int8               `json:"order" bson:"order" validate:"min=1" binding:"required" `
	TeacherId int                `json:"teacher_id" bson:"teacher_id" validate:"required" binding:"required" `
	Videos    *[]Video           `json:"videos" bson:"videos"   `
	Lectures  *[]Lecture         `json:"lectures" bson:"lectures"`
	Files     *[]Files           `json:"files" bson:"files"`
	Date
}

type Lecture struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	Name      string             `json:"name"`
	Content   string             `json:"content" bson:"content" validate:"required,min=250" binding:"required, min=250"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required" binding:"required"`
	TeacherId int                `json:"teacher_id" bson:"teacher_id" validate:"required" binding:"required"`
	IsPublic  bool               `json:"is_public" bson:"is_public" validate:"default=false" binding:"required"`
	Date
}
type Video struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Url       string             `json:"url" bson:"url" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" binding:"required"`
	TeacherId int                `json:"teacher_id" bson:"teacher_id" binding:"required"`
	Date
}
type Files struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Url       string             `json:"url" bson:"url" binding:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" binding:"required"`
	TeacherId int                `json:"teacher_id" bson:"teacher_id" binding:"required"`
	Group     int8               `json:"group" bson:"group" validate:"required"`
	Date
}

type StudentNote struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	StudentID int                `json:"student_id" bson:"student_id" binding:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" binding:"required" binding:"required"`
	Content   string             `json:"content" bson:"content" validate:"required" binding:"required"`
}
type ExtendedSection struct {
	Section

	Files    []Files        `json:"files"`
	Videos   []Video        `json:"videos"`
	Lectures []Lecture      `json:"contents"`
	Notes    *[]StudentNote `json:"note"`
}

// ---------------------- API ----------------------

type SectionResponse struct {
	Section `json:"sections"`
}
type SectionDetailsResponse struct {
	ExtendedSection `json:"section"`
}
