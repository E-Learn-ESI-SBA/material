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

func TestCreateLecture(t *testing.T) {
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       "12",
		Avatar:   "https://www.google.com",
	}
	const secret = "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"
	const url = "http://localhost:8080/lecture/"
	const id1 = "6647fc44ee41dfb53b32c9ed"
	const notValidId = "66357cea8551700b4asc5885"
	lecture := fixtures.GetLectures()[1]
	authToken, _ := utils.GenerateToken(user, secret)
	t.Run("Success: Creating  Lecture", func(t *testing.T) {
		jsonLecture, errJ := json.Marshal(lecture)
		if errJ != nil {
			t.Errorf("Error marshalling section: %v", errJ.Error())
		}
		req, errQ := http.NewRequest("POST", url+"?sectionId="+id1, bytes.NewBuffer(jsonLecture))
		if errQ != nil {
			t.Errorf("Error creating request: %v", errQ.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		resp, errR := client.Do(req)
		if errR != nil {
			t.Errorf("Error sending request: %v", errR.Error())
		}
		resBody, _ := io.ReadAll(resp.Body)
		t.Logf("Response: %v", string(resBody))
		var apiResponse interfaces.APIResponse
		err := json.Unmarshal(resBody, &apiResponse)
		if err != nil {
			t.Errorf("Error unmarshalling response: %v", err.Error())
		}
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, shared.LECTURE_CREATED, apiResponse.Message)
		// Clean up
		defer t.Cleanup(func() {
			err := resp.Body.Close()
			if err != nil {
				return
			}
		})

	})
	t.Run("Failure: Empty Section ID", func(t *testing.T) {
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			t.Errorf("Error creating request: %v", err.Error())
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Error sending request: %v", err.Error())
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Failure: Invalid Section ID", func(t *testing.T) {
		jsonLecture, errJ := json.Marshal(lecture)
		if errJ != nil {
			t.Errorf("Error marshalling lecture: %v", errJ.Error())
		}
		req, errQ := http.NewRequest("POST", url+"?sectionId="+notValidId, bytes.NewBuffer(jsonLecture))
		if errQ != nil {
			t.Errorf("Error creating request: %v", errQ.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		resp, errR := client.Do(req)
		if errR != nil {
			t.Errorf("Error sending request: %v", errR.Error())
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Failure: Invalid Request Body", func(t *testing.T) {
		req, err := http.NewRequest("POST", url+"?sectionId="+id1, bytes.NewBuffer([]byte("invalid body")))
		if err != nil {
			t.Errorf("Error creating request: %v", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Error sending request: %v", err.Error())
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
	})
}
