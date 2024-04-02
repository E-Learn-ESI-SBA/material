package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"madaurus/dev/material/app/utils"
)

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
