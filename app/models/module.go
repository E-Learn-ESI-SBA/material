package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Module struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" validate:"" `
	Year        int8               `json:"year"`
	Speciality  *string            `json:"speciality,omitempty"`
	Semester    int8               `json:"semester"`
	Coefficient int8               `json:"coefficient"`
	TeacherId   int                `json:"teacher_id" bson:"teacher_id" validate:"required"`
	Instructors *[]int             `json:"instructors,omitempty" bson:"instructors"`
	IsPublic    bool               `json:"isPublic" bson:"isPublic" validate:"default=false"`
	Image       *string            `json:"image,omitempty"`
	Date
}

type ExtendedModule struct {
	Module
	Courses *[]ExtendCourse `json:"courses,omitempty" bson:"courses" validate:"default=[]"`
}
