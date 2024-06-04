package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Mark struct {
	Type  string            `json:"type" bson:"type"`
	Attrs map[string]string `json:"attrs,omitempty" bson:"attrs,omitempty"`
}

type Node struct {
	Type    string                 `json:"type" bson:"type"`
	Attrs   map[string]interface{} `json:"attrs,omitempty" bson:"attrs,omitempty"`
	Content []Node                 `json:"content,omitempty" bson:"content,omitempty"`
	Text    string                 `json:"text,omitempty" bson:"text,omitempty"`
	Marks   []Mark                 `json:"marks,omitempty" bson:"marks,omitempty"`
}

type Content struct {
	Type    string `json:"type" bson:"type"`
	Content []Node `json:"content" bson:"content"`
}

type Lecture struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Groups    []string           `json:"groups" bson:"groups" binding:"required"`
	Name      string             `json:"name"`
	Content   Content            `json:"content" bson:"content"  binding:"required"`
	TeacherId string             `json:"teacher_id" bson:"teacher_id" validate:"required" binding:"required"`
	IsPublic  bool               `json:"is_public" bson:"is_public"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	Year      string             `json:"year" bson:"year", binding:"required"`
}

type StudentNote struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	StudentID string             `json:"student_id" bson:"student_id" binding:"required"`
	Content   Content            `json:"content" bson:"content"  binding:"required"`
}
