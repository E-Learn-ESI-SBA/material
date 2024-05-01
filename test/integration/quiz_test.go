package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/utils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Res struct {
	Message string `json:"message"`
}

var globalQuiz models.Quiz
var globalModule models.Module

// var globalTeacher utils.LightUser

var admin utils.LightUser = utils.LightUser{
	Username: "admin",
	Role:     "Admin",
	Email:    "admin@gmail.com",
	ID:       0,
}

var teacher1 utils.LightUser = utils.LightUser{
	Username: "mhammed",
	Role:     "Teacher",
	Email:    "f.mhammed@gmail.com",
	ID:       1,
}

var teacher2 utils.LightUser = utils.LightUser{
	Username: "poysa",
	Role:     "Teacher",
	Email:    "y.poysa@gmail.com",
	ID:       2,
}

var student1 utils.LightUser = utils.LightUser{
	Username: "godsword",
	Role:     "Student",
	Email:    "godsword@gmail.com",
	ID:       3,
}

var student2 utils.LightUser = utils.LightUser{
	Username: "ayoub",
	Role:     "Student",
	Email:    "ayoub@gmail.com",
	ID:       4,
}
var secretKey string = "A1B2C3D4E5F6G7H8I9J0K"
	
var err error

var adminToken string
var teacher1Token string
var teacher2Token string
var student1Token string
var student2Token string




func TestCreateQuiz(t *testing.T) {  
	adminToken, err = utils.GenerateToken(admin, secretKey)
	if err != nil {
		//throw err and test failed
		log.Printf("Error: %v\n", err)
		panic(err)
	}


	teacher1Token, err = utils.GenerateToken(teacher1, secretKey)
	if err != nil {
		//throw err and test failed
		log.Printf("Error: %v\n", err)
		panic(err)
	}

	if err != nil {
		log.Printf("Error: %v\n", err)
		panic(err)
	}

	globalModule = models.Module{
		ID: primitive.NewObjectID(),
		Name: "archi",
		Year: 3,
		Speciality: nil,
		Semester: 1,
		Coefficient: 4,
		TeacherId: teacher1.ID,
		Instructors: nil,
		IsPublic: true,
		Plan: []string{"plan1", "plan2", "plan3"},
		Image: nil,
	}


	jsonModule, _ := json.Marshal(globalModule)
	req, _ := http.NewRequest(
		"POST",
		"http://localhost:8080/modules/create",
		bytes.NewReader(jsonModule),
	)

	req.Header.Set("Authorization", "Bearer " + adminToken)
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	mockResponse := `{"message":"Module Created Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))

	globalQuiz = models.Quiz{
		ID: primitive.NewObjectID(),
		ModuleId: globalModule.ID,
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

	jsonQuiz, _ := json.Marshal(globalQuiz)
	req, _ = http.NewRequest(
		"POST",
		"http://localhost:8080/quizes/create",
		bytes.NewReader(jsonQuiz),
	)
	req.Header.Set("Authorization", "Bearer " + teacher1Token)
	//this should succeed
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	mockResponse = `{"message":"Quiz Created Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))
	// creating a quiz with a non existing teacher/module combination
	// should return an error
	quiz := globalQuiz
	quiz.ModuleId = primitive.NewObjectID()

	jsonQuiz, _ = json.Marshal(quiz)
	req, _ = http.NewRequest(
		"POST",
		"http://localhost:8080/quizes/create",
		bytes.NewReader(jsonQuiz),
	)
	req.Header.Set("Authorization", "Bearer " + teacher1Token)

	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	mockResponse = `{"error":"module not found"}`
	assert.Equal(t, mockResponse, string(responseData))
}

func TestGetQuizesByModuleId(t *testing.T) {
	
	url := "http://localhost:8080/quizes/module/" + globalModule.ID.Hex()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + teacher1Token)

	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	
	var resQuizes []models.Quiz
	json.Unmarshal(responseData, &resQuizes)
	assert.Equal(t, globalQuiz.ID, resQuizes[0].ID)
}

func TestGetQuizById(t *testing.T) {
	
	url := "http://localhost:8080/quizes/get/" + globalQuiz.ID.Hex()	
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + teacher1Token)

	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	
	var resQuiz models.Quiz
	json.Unmarshal(responseData, &resQuiz)
	assert.Equal(t, globalQuiz.ID, resQuiz.ID)
}

func TestUpdateQuiz(t *testing.T) {

	updatedQuiz := globalQuiz
	updatedQuiz.Title = "updated title..."
	updatedQuiz.Instructions = "updated instructions..."


	jsonQuiz, _ := json.Marshal(updatedQuiz)
	req, _ := http.NewRequest(
		"PUT",
		"http://localhost:8080/quizes/update",
		bytes.NewReader(jsonQuiz),
	)
	req.Header.Set("Authorization", "Bearer " + teacher1Token)

	//this should succeed
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	mockResponse := `{"message":"Quiz Updated Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))

	teacher2Token, err = utils.GenerateToken(teacher2, secretKey)
	if err != nil {
		//throw err and test failed
		log.Printf("Error: %v\n", err)
		panic(err)
	}

	//this should return an error
	req.Header.Set("Authorization", "Bearer " + teacher2Token)
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	mockResponse = `{"error":"Unauthorized"}`
	assert.Equal(t, mockResponse, string(responseData))

}

func TestSubmitQuizAnswers(t *testing.T) {
	// create two students
	student1Token, err = utils.GenerateToken(student1, secretKey)
	if err != nil {
		log.Printf("Error: %v\n", err)
		panic(err)
	}
	student2Token, err = utils.GenerateToken(student2, secretKey)
	if err != nil {
		log.Printf("Error: %v\n", err)
		panic(err)
	}

	// student1 submits the quiz with correct answers
	submission := models.Submission{
		ID: primitive.NewObjectID(),
		StudentId: student1.ID,
		QuizId: globalQuiz.ID,
		Answers: []models.Answer{
			{
				QuestionId: globalQuiz.Questions[0].ID,
				Choices: []int{0},
			},
		},
	}

	jsonSubmission, _ := json.Marshal(submission)
	url := "http://localhost:8080/quizes/" + globalQuiz.ID.Hex() + "/submit"
	req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonSubmission))
	req.Header.Set("Authorization", "Bearer " + student1Token)
	// this shoud succeed
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	mockResponse := `{"message":"Quiz Answer Submitted Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))


	// student2 submits the quiz with wrong answers
	submission = models.Submission{
		ID: primitive.NewObjectID(),
		StudentId: student2.ID,
		QuizId: globalQuiz.ID,
		Answers: []models.Answer{
			{
				QuestionId: globalQuiz.Questions[0].ID,
				Choices: []int{0, 1},
			},
		},
	}

	jsonSubmission, _ = json.Marshal(submission)
	req, _ = http.NewRequest("POST", url, bytes.NewReader(jsonSubmission))
	req.Header.Set("Authorization", "Bearer " + student2Token)
	// this shoud succeed
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	mockResponse = `{"message":"Quiz Answer Submitted Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))

	// submit again using student1
	// this should return an error
	req, _ = http.NewRequest("POST", url, bytes.NewReader(jsonSubmission))
	req.Header.Set("Authorization", "Bearer " + student1Token)
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	mockResponse = `{"error":"student already submitted the quiz answers"}`
	assert.Equal(t, mockResponse, string(responseData))
}

func TestDeleteQuiz(t *testing.T) {

	url := "http://localhost:8080/quizes/delete/" + globalQuiz.ID.Hex()
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer " + teacher2Token)
	//this should return an error
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	mockResponse := `{"error":"Unauthorized"}`
	assert.Equal(t, mockResponse, string(responseData))

	// this should succeed
	req.Header.Set("Authorization", "Bearer " + teacher1Token)
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	mockResponse = `{"message":"Quiz Deleted Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))
}