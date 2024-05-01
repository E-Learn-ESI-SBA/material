package fixtures

import (
	"madaurus/dev/material/app/models"
)

func GetCourse() []models.Course {
	var courses []models.Course
	courses = append(courses, models.Course{
		Name:        "Course 1",
		Description: "Course 1 Description",
	},
		models.Course{
			Name:        "Course 2",
			Description: "Course 2 Description",
		},

		models.Course{
			Name:        "Course 3",
			Description: "Course 3 Description",
		},
		models.Course{
			Name:        "Course 4",
			Description: "Course 4 Description",
		},
		models.Course{
			Name:        "Course 5",
			Description: "Course 5 Description",
		},
		models.Course{
			Name:        "Course 6",
			Description: "Course 6 Description",
		},
	)
	return courses
}
