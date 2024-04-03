package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	ID
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Plan        []string           `json:"plan" validate:"required,min=1" bson:"plan"`
	ModuleId    primitive.ObjectID `json:"module_id" bson:"module_id"`
	Date
}

type ExtendCourse struct {
	Course
	Sections *[]Section
}
