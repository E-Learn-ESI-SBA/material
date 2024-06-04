package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Video struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Groups    []string           `json:"groups" bson:"groups" binding:"required"`
	Url       string             `json:"url" bson:"url" validate:"required"`
	TeacherId string             `json:"teacher_id" bson:"teacher_id"`
	Name      string             `json:"name" bson:"name"  binding:"required"`
	Score     int32              `json:"score" bson:"score" binding:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
}
