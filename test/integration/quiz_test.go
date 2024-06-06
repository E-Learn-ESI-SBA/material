package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"madaurus/dev/material/test/fixtures"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


type Res struct {
	Message string `json:"message"`
}

var globalQuiz models.Quiz
var globalModule models.Module


var admin utils.LightUser = fixtures.GetAdmins()[0]

var teacher1 utils.LightUser = fixtures.GetTeachers()[1]
var teacher2 utils.LightUser = fixtures.GetTeachers()[0]
var student1 utils.LightUser = fixtures.GetStudents()[0]
var student2 utils.LightUser = fixtures.GetStudents()[1]
var secretKey string = "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"
	
var adminToken string
var teacher1Token string
var teacher2Token string
var student1Token string
var student2Token string




func TestCreateQuiz(t *testing.T) {  

	var err error
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
	

	globalModule = fixtures.GetModules()[1]

	jsonModule, _ := json.Marshal(globalModule)
	req, _ := http.NewRequest(
		"POST",
		"http://localhost:8080/modules",
		bytes.NewReader(jsonModule),
	)

	req.Header.Set("Authorization", "Bearer " + adminToken)
	req.Header.Set("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)

	var apiResponse interfaces.APiSuccess
	err = json.Unmarshal(responseData, &apiResponse)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	assert.Equal(t, shared.CREATE_MODULE, apiResponse.Message)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// fetch the created module 
	// globalModule.TeacherId = teacher1.ID
	// assert.Equal(t, globalModule.TeacherId, teacher1.ID)
	// req, _ = http.NewRequest(
	// 	"GET",
	// 	"http://localhost:8080/modules/teacher",
	// 	nil,
	// )
	// req.Header.Set("Authorization", "Bearer " + teacher1Token)
	// res, _ = http.DefaultClient.Do(req)
	// responseData, _ = io.ReadAll(res.Body)
	// // unmarshal the response
	// log.Printf("Response: %v\n", string(responseData))
	// var resModules []models.Module
	// err = json.Unmarshal(responseData, &resModules)
	// if err != nil {
	// 	t.Errorf("Error unmarshalling response body: %v", err)
	// }
	// // globalModule.ID = resModules[0].ID
	t.Run("Create Quiz", func(t *testing.T) {
		moduleHexId := globalModule.ID.Hex()
		globalQuiz = fixtures.GetQuiz(moduleHexId)
		jsonQuiz, _ := json.Marshal(globalQuiz)
		req, _ = http.NewRequest(
			"POST",
			"http://localhost:8080/quizes",
			bytes.NewReader(jsonQuiz),
		)
		req.Header.Set("Authorization", "Bearer " + teacher1Token)
		//this should succeed
		res, _ = http.DefaultClient.Do(req)
		responseData, _ = io.ReadAll(res.Body)

		var apiResponse interfaces.APiSuccess
		err = json.Unmarshal(responseData, &apiResponse)
		if err != nil {
			t.Errorf("Error unmarshalling response body: %v", err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, shared.QUIZ_CREATED, apiResponse.Message)
	})
	// // creating a quiz with a non existing teacher/module combination
	// // should return an error
	// quiz := globalQuiz
	// quiz.ModuleId = primitive.NewObjectID()

	// jsonQuiz, _ = json.Marshal(quiz)
	// req, _ = http.NewRequest(
	// 	"POST",
	// 	"http://localhost:8080/quizes/create",
	// 	bytes.NewReader(jsonQuiz),
	// )
	// req.Header.Set("Authorization", "Bearer " + teacher1Token)

	// res, _ = http.DefaultClient.Do(req)
	// responseData, _ = io.ReadAll(res.Body)
	// mockResponse = `{"error":"module not found"}`
	// assert.Equal(t, mockResponse, string(responseData))
}

// func TestGetQuizesByModuleId(t *testing.T) {
	
// 	url := "http://localhost:8080/quizes/module/" + globalModule.ID.Hex()
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Authorization", "Bearer " + teacher1Token)

// 	res, _ := http.DefaultClient.Do(req)
// 	responseData, _ := io.ReadAll(res.Body)
	
// 	var resQuizes []models.Quiz
// 	json.Unmarshal(responseData, &resQuizes)
// 	assert.Equal(t, globalQuiz.ID, resQuizes[0].ID)
// }

func TestGetQuizesByAdmin(t *testing.T) {
	url := "http://localhost:8080/quizes/admin"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + adminToken)

	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)

	var resQuizes []models.Quiz
	json.Unmarshal(responseData, &resQuizes)
	globalQuiz.ID = resQuizes[0].ID
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetQuizById(t *testing.T) {
	
	url := "http://localhost:8080/quizes/" + globalQuiz.ID.Hex()	
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + teacher1Token)

	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	

	type Res struct {
		passed bool
		quiz models.Quiz
	}

	var resQuiz Res
	json.Unmarshal(responseData, &resQuiz)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	// assert.Equal(t, globalQuiz.ID, quiz.ID)
}

func TestGetQuizQuestions(t *testing.T) {
	log.Println(globalQuiz.ID.Hex())
	url := "http://localhost:8080/quizes/" + globalQuiz.ID.Hex() + "/questions"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + teacher1Token)
	
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)

	log.Printf("Response: %v\n", string(responseData))

	var questionsRes struct {
		Questions []models.Question
		Duration int
		Title string
	}

	err := json.Unmarshal(responseData, &questionsRes)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, globalQuiz.Questions[0].Body, questionsRes.Questions[0].Body)
}

func TestUpdateQuiz(t *testing.T) {

	var quizUpdates models.QuizUpdate
	quizUpdates.Title = "updated title..."
	quizUpdates.Instructions = "updated instructions..."
	quizUpdates.StartDate = globalQuiz.StartDate
	quizUpdates.EndDate = globalQuiz.EndDate
	quizUpdates.Duration = globalQuiz.Duration


	jsonQuiz, _ := json.Marshal(quizUpdates)
	req, _ := http.NewRequest(
		"PUT",
		"http://localhost:8080/quizes/" + globalQuiz.ID.Hex(),
		bytes.NewReader(jsonQuiz),
	)
	req.Header.Set("Authorization", "Bearer " + teacher1Token)

	//this should succeed
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	var resRes interfaces.APiSuccess
	err := json.Unmarshal(responseData, &resRes)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, shared.QUIZ_UPDATED, resRes.Message)

	teacher2Token, err = utils.GenerateToken(teacher2, secretKey)
	if err != nil {
		//throw err and test failed
		log.Printf("Error: %v\n", err)
		panic(err)
	}

	// this should return an error
	req.Header.Set("Authorization", "Bearer " + teacher2Token)
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)

	var resErr interfaces.APiError
	err = json.Unmarshal(responseData, &resErr)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}


	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, shared.UNAUTHORIZED, resErr.Error)

}

func TestSubmitQuizAnswers(t *testing.T) {
	// create two students
	var err error
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
		StudentId: student1.ID,
		QuizId: globalQuiz.ID,
		Answers: []models.Answer{
			{
				QuestionId: globalQuiz.Questions[0].ID,
				Choices: []string{"0"},
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

	log.Printf("Response: %v\n", string(responseData))

	var resRes interfaces.APiSuccess
	err = json.Unmarshal(responseData, &resRes)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, shared.QUIZ_ANSWER_SUBMITTED, resRes.Message)


	// student2 submits the quiz with wrong answers
	submission = models.Submission{
		StudentId: student2.ID,
		QuizId: globalQuiz.ID,
		Answers: []models.Answer{
			{
				QuestionId: globalQuiz.Questions[0].ID,
				Choices: []string{"0", "1"},
			},
		},
	}

	jsonSubmission, _ = json.Marshal(submission)
	req, _ = http.NewRequest("POST", url, bytes.NewReader(jsonSubmission))
	req.Header.Set("Authorization", "Bearer " + student2Token)
	// this shoud succeed
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)

	log.Printf("Response: %v\n", string(responseData))

	err = json.Unmarshal(responseData, &resRes)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}	
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, shared.QUIZ_ANSWER_SUBMITTED, resRes.Message)

	// submit again using student1
	// this should return an error
	req, _ = http.NewRequest("POST", url, bytes.NewReader(jsonSubmission))
	req.Header.Set("Authorization", "Bearer " + student1Token)
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	
	log.Printf("Response: %v\n", string(responseData))

	var resErr interfaces.APiError
	err = json.Unmarshal(responseData, &resErr)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, shared.QUIZ_ANSWER_ALREADY_SUBMITTED, resErr.Error)
}

func TestGetQuizResults(t *testing.T) {
	url := "http://localhost:8080/quizes/" + globalQuiz.ID.Hex() + "/teacher"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + teacher1Token)
	// this should succeed
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	var resResults []models.Submission
	json.Unmarshal(responseData, &resResults)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, len(resResults))

	// same request with teacher2 token
	// this should return an error
	teacher2Token, err := utils.GenerateToken(teacher2, secretKey)
	if err != nil {
		log.Printf("Error: %v\n", err)
		panic(err)
	}
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + teacher2Token)
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	var resErr interfaces.APiError
	err = json.Unmarshal(responseData, &resErr)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, shared.QUIZ_NOT_FOUND, resErr.Error)
}

func TestGetQuizResultByStudentId(t *testing.T) {
	url := "http://localhost:8080/quizes/" + globalQuiz.ID.Hex() + "/student"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + student1Token)
	// this should succeed
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	var resSubmission models.Submission
	json.Unmarshal(responseData, &resSubmission)
	// this should fail bcz end date is not yet reached
	// assert.Equal(t, student1.ID, resSubmission.StudentId)
	var resErr interfaces.APiError
	err := json.Unmarshal(responseData, &resErr)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, shared.QUIZ_STILL_ONGOING, resErr.Error)

	//sleep for 2 seconds
	time.Sleep(2 * time.Second)
	// this should succeed
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	json.Unmarshal(responseData, &resSubmission)
	// assert.Equal(t, student1.ID, resSubmission.StudentId)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// same request with teacher2 token
	// this should return an error
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + teacher2Token)
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	
	err = json.Unmarshal(responseData, &resErr)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, shared.QUIZ_ANSWER_NOT_FOUND, resErr.Error)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetQuizesResultByStudentId(t *testing.T) {
	url := "http://localhost:8080/quizes/student"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + student1Token)

	// this should succeed
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	var resResults []models.Submission
	err := json.Unmarshal(responseData, &resResults)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, student1.ID, resResults[0].StudentId)
}

func TestDeleteQuiz(t *testing.T) {

	url := "http://localhost:8080/quizes/" + globalQuiz.ID.Hex()
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer " + teacher2Token)
	//this should return an error
	res, _ := http.DefaultClient.Do(req)
	responseData, _ := io.ReadAll(res.Body)
	var resErr interfaces.APiError
	err := json.Unmarshal(responseData, &resErr)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, shared.UNAUTHORIZED, resErr.Error)

	// this should succeed
	req.Header.Set("Authorization", "Bearer " + teacher1Token)
	res, _ = http.DefaultClient.Do(req)
	responseData, _ = io.ReadAll(res.Body)
	
	var resRes interfaces.APiSuccess
	err = json.Unmarshal(responseData, &resRes)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, shared.QUIZ_DELETED, resRes.Message)
}