package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Files struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Url  string             `json:"url" bson:"url"`
	Name string             `json:"name" bson:"name"  binding:"required"`
	//	SectionId primitive.ObjectID `json:"section_id" bson:"section_id"`
	Type      string    `json:"type" bson:"type"`
	TeacherId string    `json:"teacher_id" bson:"teacher_id"`
	Groups    []string  `json:"groups" bson:"groups" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
