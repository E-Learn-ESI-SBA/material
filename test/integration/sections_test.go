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
	}
	const secret = "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"
	const url = "http://localhost:8080/section/"
	const id1 = "6632042f17c6d49ef81973f9"
	const id2 = "66327ca0547f7f30e5b2a733"
	section := fixtures.GetSections()[0]
	authToken, _ := utils.GenerateToken(user, secret)

	t.Run("Success: Creating  Section", func(t *testing.T) {
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