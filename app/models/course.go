package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	Name        string             `json:"name"`
	Description string             `json:"description"`
	ModuleId    primitive.ObjectID `json:"module_id" bson:"module_id"`
	Date
}

type ExtendCourse struct {
	Course
	Sections *[]Section
}
