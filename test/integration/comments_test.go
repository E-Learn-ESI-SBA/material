package integration

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"madaurus/dev/material/test/fixtures"
	"net/http"
	"testing"
)

func TestCreateComment(t *testing.T) {
	courseID := "66357c958551700b413c587c" // Replace with a valid course ID
	comment := fixtures.GetComment()[0]
	jsonComment, err := json.Marshal(comment)
	if err != nil {
		t.Errorf("Error marshaling comment: %v", err)
	}
	user := utils.LightUser{
		Email:    "moha@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       "2",
		Avatar:   "https://www.google.com",
		Year:     "2021",
		Group:    "1A",
	}
	const secret = "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"
	authToken, _ := utils.GenerateToken(user, secret)
	t.Run("Success", func(t *testing.T) {
		url := "http://localhost:8080/comments/?courseId=" + courseID
		req, err := http.NewRequest("POST", url, bytes.NewReader(jsonComment))
		if err != nil {
			t.Errorf("Error creating request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}

		var apiResponse interfaces.APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResponse)
		if err != nil {
			t.Errorf("Error decoding response: %v", err)
		}

		assert.Equal(t, shared.COMMENT_CREATED, apiResponse.Message)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid Course ID", func(t *testing.T) {
		url := "http://localhost:8080/comments/?courseId=invalid"
		req, err := http.NewRequest("POST", url, bytes.NewReader(jsonComment))
		if err != nil {
			t.Errorf("Error creating request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}

		var apiResponse interfaces.APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResponse)
		if err != nil {
			t.Errorf("Error decoding response: %v", err)
		}

		assert.Equal(t, shared.INVALID_ID, apiResponse.Message)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invalid Comment Body", func(t *testing.T) {
		url := "http://localhost:8080/comments/?courseId=" + courseID
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			t.Errorf("Error creating request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		var apiResponse interfaces.APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResponse)
		if err != nil {
			t.Errorf("Error decoding response: %v", err)
		}

		assert.Equal(t, shared.INVALID_BODY, apiResponse.Message)
		//assert.Equal(t, http.StatusNotAcceptable, resp.Status)
	})
}
