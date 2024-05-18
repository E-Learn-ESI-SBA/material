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

func TestCreateSection(t *testing.T) {
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       "12",
		Avatar:   "https://www.google.com",
		Year:     "2021",
		Group:    "1A",
	}
	const secret = "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"
	const url = "http://localhost:8080/section/"
	const id1 = "6647fbedee41dfb53b32c9e6"
	const id2 = "6647fc00ee41dfb53b32c9e7"
	section := fixtures.GetSections()[0]
	authToken, _ := utils.GenerateToken(user, secret)

	t.Run("Success: Creating  Section", func(t *testing.T) {
		jsonSection, errJ := json.Marshal(section)
		if errJ != nil {
			t.Errorf("Error marshalling section: %v", errJ.Error())
		}

		req, errQ := http.NewRequest("POST", url+"?courseId="+id2, bytes.NewBuffer(jsonSection))

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
		assert.Equal(t, shared.CREATE_SECTION, apiResponse.Message)

		// Clean up
		defer t.Cleanup(func() {
			err := resp.Body.Close()
			if err != nil {
				return
			}
		})

	})

	t.Run("Failure: Creating Section with invalid courseId", func(t *testing.T) {
		// Create a section with invalid courseId
		jsonSection, errJ := json.Marshal(section)
		if errJ != nil {
			t.Errorf("Error marshalling section: %v", errJ.Error())

		}
		req, errQ := http.NewRequest("POST", url+"?courseId=", bytes.NewBuffer(jsonSection))
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
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, shared.INVALID_ID, apiResponse.Message)

	})

	t.Run("Failure: Creating Section with invalid body", func(t *testing.T) {
		// Create a section with invalid body
		section.Name = ""
		jsonSection, errJ := json.Marshal(section)
		if errJ != nil {
			t.Errorf("Error marshalling section: %v", errJ.Error())
		}
		req, errQ := http.NewRequest("POST", url+"?courseId="+id1, bytes.NewBuffer(jsonSection))
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
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, shared.INVALID_BODY, apiResponse.Message)

	})
}
