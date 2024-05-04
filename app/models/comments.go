package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"madaurus/dev/material/app/utils"
	"time"
)

type Comments struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Content   string             `json:"content" bson:"content" validate:"required" binding:"required, min=250"`
	UserId    string             `json:"user_id" bson:"user_id" validate:"required" binding:="required"`
	CourseId  primitive.ObjectID `json:"course_id" bson:"course_id" validate:"required" binding:"required"`
	IsEdited  bool               `json:"is_edited" bson:"is_edited" validate:"default=false" binding:"default=false"`
	User      utils.LightUser    `json:"user" bson:"user"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	Replays   *[]Reply           `json:"replays" bson:"replays"`
}

type Reply struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Content  string             `json:"content" bson:"content" binding:"required"`
	UserId   string             `json:"user_id" bson:"user_id" binding:"required"`
	User     utils.LightUser    `json:"user" bson:"user"`
	IsEdited bool               `json:"is_edited" bson:"is_edited" binding:"default=false"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
