package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	Name        string             `json:"name" binding:"required" bson:"name" validate:"required"`
	Description string             `json:"description" binding:"required" bson:"description" validate:"required"`
	ModuleId    primitive.ObjectID `json:"module_id" bson:"module_id"`
	Date
}

type ExtendCourse struct {
	Course
	Sections *[]Section
}

type UltraCourse struct {
	Course
	Sections *[]ExtendedSection `json:"sections"`
}
