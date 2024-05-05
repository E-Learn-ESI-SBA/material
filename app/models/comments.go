package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"madaurus/dev/material/app/utils"
	"time"
)

type Comments struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Content   string             `json:"content" bson:"content"  binding:"required"`
	UserId    string             `json:"user_id" bson:"user_id" `
	CourseId  primitive.ObjectID `json:"course_id" bson:"course_id"  binding:"required"`
	IsEdited  bool               `json:"is_edited" bson:"is_edited"`
	User      utils.LightUser    `json:"user" bson:"user"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	Replays   []Reply            `json:"replays" bson:"replays"`
}

type Reply struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Content  string             `json:"content" bson:"content" binding:"required"`
	UserId   string             `json:"user_id" bson:"user_id"`
	User     utils.LightUser    `json:"user" bson:"user"`
	IsEdited bool               `json:"is_edited" bson:"is_edited" `

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
