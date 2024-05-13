package fixtures

import (
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


 

func GetTeachers() []utils.LightUser {
	return []utils.LightUser{
		utils.LightUser{
			Username: "mhammed",
			Role:     "Teacher",
			Email:    "f.mhammed@gmail.com",
			ID:       1,
		},
		utils.LightUser{
			Username: "poysa",
			Role:     "Teacher",
			Email:    "y.poysa@gmail.com",
			ID:       2,
		},
	}

}


func GetStudents() []utils.LightUser {
	return []utils.LightUser{
		utils.LightUser{
			Username: "godsword",
			Role:     "Student",
			Email:    "godsword@gmail.com",
			ID:       3,
		},
		utils.LightUser{
			Username: "ayoub",
			Role:     "Student",
			Email:    "ayoub@gmail.com",
			ID:       4,
		},
	}
}


func GetAdmins() []utils.LightUser {
	return []utils.LightUser{
		{
			Username: "admin",
			Role:     "Admin",
			Email:    "admin@gmail.com",
			ID:       0,
		},
	}
}


func GetQuiz(moduleId primitive.ObjectID) models.Quiz {
	return models.Quiz{
		ID: primitive.NewObjectID(),
		ModuleId: moduleId,
		Title: "quiz_goes_brr",
    	Instructions: "some instructions...",
    	QuestionCount: 20,
		MaxScore: 100,
    	StartDate: time.Now(),
    	EndDate: time.Now().Add(time.Hour * 1), // after one hour
    	Duration: 100,
		Questions: []models.Question{
			{
				ID: primitive.NewObjectID(),
				Body: "what is the capital of france?",
				Description: "extra info (optional)",
				Score: 100,
				Options: []string{"paris", "london", "berlin", "madrid"},
				CorrectIdxs: []int{0},
			},
		},
		Grades: []models.Grade{
			{
				Min:    0,
				Max:    33,
				Grade:  "C",
				IsPass: false,
			},
			{
				Min:    34,
				Max:    66,
				Grade:  "B",
				IsPass: true,
			},
			{
				Min:    67,
				Max:    100,
				Grade:  "A",
				IsPass: true,
			},
		},
	}
}