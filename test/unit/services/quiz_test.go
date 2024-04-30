package services_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"madaurus/dev/material/app/models"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Res struct {
	Message string `json:"message"`
}

func CreateQuiz() {
	quiz := models.Quiz{
		ModuleId: primitive.NewObjectID(),
    	TeacherId: 3,
		Title: "quiz_goes_brr",
    	Instructions: "some instructions...",
    	MinScore: 50,
    	QuestionCount: 20,
    	StartDate: time.Now(),
    	EndDate: time.Now(),
    	Duration: 100,
	}

	// encode to json
	jsonQuiz, _ := json.Marshal(quiz)
	res, _ := http.Post("http://localhost:8080/quizes/create", "application/json", bytes.NewReader(jsonQuiz))

	responseData, _ := ioutil.ReadAll(res.Body)
	var response Res
	json.Unmarshal(responseData, &response)
	log.Printf("Response: %v\n", response)
}

func TestCreateQuiz(t *testing.T) {  
	quiz := models.Quiz{
		ModuleId: primitive.NewObjectID(),
    	TeacherId: 3,
		Title: "quiz_goes_brr",
    	Instructions: "some instructions...",
    	MinScore: 50,
    	QuestionCount: 20,
    	StartDate: time.Now(),
    	EndDate: time.Now(),
    	Duration: 100,
	}
	// encode to json
	jsonQuiz, _ := json.Marshal(quiz)
	res, _ := http.Post("http://localhost:8080/quizes/create", "application/json", bytes.NewReader(jsonQuiz))

	responseData, _ := ioutil.ReadAll(res.Body)
	var response Res
	json.Unmarshal(responseData, &response)
	mockResponse := `{"message":"Quiz Created Successfully"}`
	assert.Equal(t, mockResponse, string(responseData))
}

func TestGetQuizById(t *testing.T) {
	// fetch all quizes
	res, _ := http.Get("http://localhost:8080/quizes/admin")
	responseData, _ := ioutil.ReadAll(res.Body)
	var quizes []models.Quiz
	json.Unmarshal(responseData, &quizes)

	// delete the first quiz
	quiz := quizes[0]

	url := "http://localhost:8080/quizes/get/" + quiz.ID.Hex()
	log.Println(url)
	res, _ = http.Get(url)
	responseData, _ = ioutil.ReadAll(res.Body)
	var quiz1 models.Quiz
	json.Unmarshal(responseData, &quiz1)
	assert.Equal(t, quiz.ID, quiz1.ID)
	// assert.NotEqual(t, mockResponse, string(responseData))
}

func TestDeleteQuiz(t *testing.T) {

}