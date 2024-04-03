package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Date struct {
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type ID struct {
	ID *primitive.ObjectID `json:"id" bson:"_id"`
}
