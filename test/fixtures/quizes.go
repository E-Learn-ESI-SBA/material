package fixtures

import (
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared/iam"
	"madaurus/dev/material/app/utils"
	"time"
)


 

func GetTeachers() []utils.LightUser {
	return []utils.LightUser{
		utils.LightUser{
			Username: "mhammed",
			Role:     iam.ROLETeacherKey,
			Email:    "f.mhammed@gmail.com",
			ID:       "2",
		},
		utils.LightUser{
			Username: "poysa",
			Role:     iam.ROLETeacherKey,
			Email:    "y.poysa@gmail.com",
			ID:       "3",
		},
	}

}


func GetStudents() []utils.LightUser {
	return []utils.LightUser{
		utils.LightUser{
			Username: "godsword",
			Role:     iam.ROLEStudentKey,
			Email:    "godsword@gmail.com",
			ID:       "3",
		},
		utils.LightUser{
			Username: "ayoub",
			Role:     iam.ROLEStudentKey,
			Email:    "ayoub@gmail.com",
			ID:       "4",
		},
	}
}


func GetAdmins() []utils.LightUser {
	return []utils.LightUser{
		{
			Username: "admin",
			Role:     iam.ROLEAdminKey,
			Email:    "admin@gmail.com",
			ID:       "2",
		},
	}
}


func GetQuiz(moduleId string) models.Quiz {
	return models.Quiz{
		ModuleId: moduleId,
		Title: "quiz_goes_brr",
    	Instructions: "some instructions...",
    	QuestionCount: 20,
		MaxScore: 100,
		MinScore: 50,
    	StartDate: time.Now(),
    	EndDate: time.Now().Add(time.Second * 2), // after two seconds
    	Duration: 100,
		Year: "some year here",
		Questions: []models.Question{
			{
				ID: "0",
				Body: "what is the capital of france?",
				Score: 100,
				Options: []models.Option{
					models.Option{
						ID: "0",
						Option: "Paris",
					},
					models.Option{
						ID: "1",
						Option: "London",
					},
					models.Option{
						ID: "2",
						Option: "Berlin",
					},
					models.Option{
						ID: "3",
						Option: "Madrid",
					},
				},
				CorrectIdxs: []string{"0"},
			},
		},
		Grades: []models.Grade{
			{
				Min:    0,
				Max:    33,
				Grade:  "C",
			},
			{
				Min:    34,
				Max:    66,
				Grade:  "B",
			},
			{
				Min:    67,
				Max:    100,
				Grade:  "A",
			},
		},
	}
}