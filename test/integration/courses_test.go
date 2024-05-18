package integration

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"madaurus/dev/material/test/fixtures"
	"net/http"
	"testing"
)

func TestCreateChapter(t *testing.T) {
	const url = "http://localhost:8080/courses"
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       "12",
		Avatar:   "sqkddkslqdns",
		Year:     "2021",
		Group:    "1A",
	}
	const module_id = "6647fb86ee41dfb53b32c9e5"
	const secret = "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"
	authToken, _ := utils.GenerateToken(user, secret)
	course := fixtures.GetCourse()[0]
	t.Run("Create Course", func(t *testing.T) {
		jsonCourse, err := json.Marshal(course)
		if err != nil {
			t.Errorf("Error in marshaling the course")

		}
		req, errQ := http.NewRequest("POST", url+"?module="+module_id, bytes.NewBuffer(jsonCourse))
		if errQ != nil {
			t.Errorf("Error in creating the request")
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, errR := client.Do(req)
		if errR != nil {
			t.Errorf("Error in sending the request")

		}
		resBody, _ := io.ReadAll(res.Body)

		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})
		var apiResponse interfaces.APIResponse
		err = json.Unmarshal(resBody, &apiResponse)
		if err != nil {
			t.Errorf("Error unmarshalling response body: %v", err)
		}
		assert.Equal(t, http.StatusCreated, res.StatusCode, "Course Created Successfully")
		assert.Equal(t, shared.CREATE_COURSE, apiResponse.Message, "Course Created Successfully With Message :")
	})
	t.Run("Create Course with missing fields", func(t *testing.T) {
		course.Name = ""
		course.Description = ""
		jsonCourse, err := json.Marshal(course)
		if err != nil {
			t.Errorf("Error in marshaling the course")

		}
		req, errQ := http.NewRequest("POST", url+"?module="+module_id, bytes.NewBuffer(jsonCourse))
		if errQ != nil {
			t.Errorf("Error in creating the request")
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, errR := client.Do(req)
		if errR != nil {
			t.Errorf("Error in sending the request")

		}
		resBody, _ := io.ReadAll(res.Body)

		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})
		var apiResponse interfaces.APIResponse
		err = json.Unmarshal(resBody, &apiResponse)
		if err != nil {
			t.Errorf("Error unmarshalling response body: %v", err)
		}
		assert.Equal(t, http.StatusNotAcceptable, res.StatusCode)
		assert.Equal(t, shared.INVALID_BODY, apiResponse.Message)
	})
	t.Run("Create Course with invalid token", func(t *testing.T) {
		jsonCourse, err := json.Marshal(course)
		if err != nil {
			t.Errorf("Error in marshaling the course")

		}
		req, errQ := http.NewRequest("POST", url+"?module="+module_id, bytes.NewBuffer(jsonCourse))
		if errQ != nil {
			t.Errorf("Error in creating the request")
		}
		req.Header.Set("Authorization", "Bearer invalid")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, errR := client.Do(req)
		if errR != nil {
			t.Errorf("Error in sending the request")

		}
		resBody, _ := io.ReadAll(res.Body)

		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})
		var apiResponse interfaces.APIResponse
		err = json.Unmarshal(resBody, &apiResponse)
		if err != nil {
			t.Errorf("Error unmarshalling response body: %v", err)
		}
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Equal(t, shared.UNAUTHORIZED, apiResponse.Message)
	})

}
