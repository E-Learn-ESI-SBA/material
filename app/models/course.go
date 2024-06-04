package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Course struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	Name        string `json:"name" binding:"required" bson:"name" validate:"required"`
	Description string `json:"description" binding:"required" bson:"description" validate:"required"`
	//  ModuleId    primitive.ObjectID `json:"module_id" bson:"module_id"`
	Sections  []Section `json:"sections" bson:"sections"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	Year      string    `json:"year" bson:"year"`
}

type ExtendCourse struct {
	Course
	Sections *[]Section
}

type UltraCourse struct {
	Course
	Sections *[]ExtendedSection `json:"sections"`
}

type RCourse struct {
	Course
	ID       string     `json:"id"`
	Sections []RSection `json:"sections"`
}

func (c *RCourse) Extract(course Course) {
	c = &RCourse{
		Course:   course,
		ID:       course.ID.Hex(),
		Sections: []RSection{},
	}
	rs := RSection{}
	for _, section := range course.Sections {
		rs.Extract(section)
		c.Sections = append(c.Sections, rs)

	}
}
