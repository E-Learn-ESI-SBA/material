package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Module struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	Name        string    `json:"name" binding:"required" bson:"name" validate:"required"`
	Year        int8      `json:"year" binding:"required" bson:"year" validate:"required"`
	Speciality  *string   `json:"speciality,omitempty"`
	Semester    int8      `json:"semester" binding:"required" bson:"semester" validate:"required"`
	Coefficient int8      `json:"coefficient" binding:"required" bson:"coefficient" validate:"required"`
	TeacherId   string    `json:"teacher_id" bson:"teacher_id" validate:"required" binding:"required" `
	Instructors *[]string `json:"instructors,omitempty" bson:"instructors" `
	IsPublic    bool      `json:"isPublic" bson:"isPublic"`
	Plan        []string  `json:"plan" binding:"required,min=1" bson:"plan"`
	Image       *string   `json:"image,omitempty"`
	Courses     []Course  `json:"courses" bson:"courses"`
	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type ExtendedModule struct {
	Module
	Courses *[]ExtendCourse `json:"courses,omitempty" bson:"courses" validate:"default=[]"`
}
