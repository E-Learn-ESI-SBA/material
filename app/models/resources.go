package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Section struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" validate:"required" binding:"required" bson:"name" `
	TeacherId string             `json:"teacher_id" bson:"teacher_id"`
	Videos    []Video            `json:"videos" bson:"videos"`
	Lectures  []Lecture          `json:"lectures" bson:"lectures"`
	Files     []Files            `json:"files" bson:"files"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
}

type Lecture struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	Name    string `json:"name"`
	Content string `json:"content" bson:"content" validate:"required,min=250" binding:"required, min=250"`
	//	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required" binding:"required"`
	TeacherId string    `json:"teacher_id" bson:"teacher_id" validate:"required" binding:"required"`
	IsPublic  bool      `json:"is_public" bson:"is_public" validate:"default=false" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
type Video struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Url       string             `json:"url" bson:"url" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" binding:"required"`
	TeacherId string             `json:"teacher_id" bson:"teacher_id" binding:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
}
type Files struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Url  string             `json:"url" bson:"url"`
	Name string             `json:"name" bson:"name" validate:"required" binding:"required"`
	//	SectionId primitive.ObjectID `json:"section_id" bson:"section_id"`
	Type      string    `json:"type" bson:"type"  binding:"default=text"`
	TeacherId string    `json:"teacher_id" bson:"teacher_id"`
	Group     string    `json:"group" bson:"group"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type StudentNote struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	StudentID string             `json:"student_id" bson:"student_id" binding:"required"`
	//	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" binding:"required" binding:"required"`
	Content string `json:"content" bson:"content" validate:"required" binding:"required"`
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
