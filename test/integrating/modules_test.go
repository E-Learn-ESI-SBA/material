package integrating

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/utils"
	"madaurus/dev/material/test/fixtures"
	"net/http"
	"testing"
)

func TestCreateModule(t *testing.T) {
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       12,
	}
	const url = "http://localhost:8080"
	const secret = "A1B2C3D4E5F6G7H8I9J0K"
	authToken, _ := utils.GenerateToken(user, secret)
	module := fixtures.GetModules()[0]
	t.Run("Success Creating Module", func(t *testing.T) {
		jsonModule, err := json.Marshal(module)
		if err != nil {
			t.Errorf("UnExpected Error: %v", err)
		}
		req, errR := http.NewRequest("POST", url+"/modules/create", bytes.NewReader(jsonModule))
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		resBody, _ := io.ReadAll(res.Body)

		// Fixed Unmarshalling
		var apiResponse interfaces.APiSuccess
		err = json.Unmarshal(resBody, &apiResponse)
		if err != nil {
			t.Errorf("Error unmarshalling response body: %v", err)
		}

		// Optional Status Code Assertion
		assert.Equal(t, 200, res.StatusCode)

		// Message Assertion
		assert.Equal(t, "Module Created Successfully", apiResponse.Message)

		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})
	})
	t.Run("Unauthorized Access (Missing Auth Token)", func(t *testing.T) {
		jsonModule, err := json.Marshal(module)
		if err != nil {
			t.Errorf("UnExpected Error: %v", err)
		}
		req, errR := http.NewRequest("POST", url+"/modules/create", bytes.NewReader(jsonModule))
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})
	t.Run("Bad Request , inavalid Body", func(t *testing.T) {
		var emptyModule models.Module
		jsonModule, err := json.Marshal(emptyModule)
		if err != nil {
			t.Errorf("UnExpected Error: %v", err)
		}
		req, errR := http.NewRequest("POST", url+"/modules/create", bytes.NewReader(jsonModule))
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		resBody, _ := io.ReadAll(res.Body)

		var apiResponse interfaces.APiSuccess
		err = json.Unmarshal(resBody, &apiResponse)
		if err != nil {
			t.Errorf("Error unmarshalling response body: %v", err)
		}

		assert.Equal(t, 400, res.StatusCode)
		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})
	})
}

// Test The Edit Module

func TestUpdateModule(t *testing.T) {
	var module models.Module
	const id = "663183f584a7d58b141442ac"

}
