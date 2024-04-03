package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"madaurus/dev/material/app/utils"
)

type Comments struct {
	ID
	Content  string             `json:"content" bson:"content" validate:"required"`
	UserId   int                `json:"user_id" bson:"user_id" validate:"required"`
	CourseId primitive.ObjectID `json:"course_id" bson:"course_id" validate:"required"`
	IsEdited bool               `json:"is_edited" bson:"is_edited" validate:"default=false"`

	Date
	Replays *[]Reply `json:"replays" bson:"replays"`
}

type Reply struct {
	ID
	Content  string          `json:"content" bson:"content" validate:"required"`
	UserId   string          `json:"user_id" bson:"user_id" validate:"required"`
	User     utils.LightUser `json:"user" bson:"user"`
	IsEdited bool            `json:"is_edited" bson:"is_edited" validate:"default=false"`

	Date
}
