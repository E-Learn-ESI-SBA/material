package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comments struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Content   string             `json:"content" bson:"content"  binding:"required,lt=500"`
	UserId    string             `json:"user_id" bson:"user_id"`
	CourseId  primitive.ObjectID `json:"course_id" bson:"course_id"`
	IsEdited  bool               `json:"is_edited" bson:"is_edited"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	Replays   []Reply            `json:"replays" bson:"replays"`
	User      User               `json:"user" bson:"user"`
}

type Reply struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Content   string             `json:"content" bson:"content" binding:"required,lt=500"`
	UserId    string             `json:"user_id" bson:"user_id"`
	IsEdited  bool               `json:"is_edited" bson:"is_edited" `
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	User      User               `json:"user" bson:"user"`
}
