package integrating

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/utils"
	"madaurus/dev/material/test/fixtures"
	"net/http"
	"testing"
	"time"
)

const url = "http://localhost:8080"

func TestCreateModule(t *testing.T) {
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       12,
	}
	const secret = "A1B2C3D4E5F6G7H8I9J0K"
	authToken, _ := utils.GenerateToken(user, secret)
	module := fixtures.GetModules()[0]
	t.Run("Success Creating Module", func(t *testing.T) {
		jsonModule, err := json.Marshal(module)
		if err != nil {
			t.Errorf("UnExpected Error: %v", err)
		}
		req, errR := http.NewRequest("POST", url+"/modules", bytes.NewReader(jsonModule))
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

		log.Printf("Response message   %v ", apiResponse.Message)

		// Optional Status Code Assertion
		assert.Equal(t, 201, res.StatusCode)

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
		req, errR := http.NewRequest("POST", url+"/modules", bytes.NewReader(jsonModule))
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
	t.Run("Bad Request , invalid Body", func(t *testing.T) {

		var emptyModule models.Module
		jsonModule, err := json.Marshal(emptyModule)
		if err != nil {
			t.Errorf("UnExpected Error: %v", err)
		}
		req, errR := http.NewRequest("POST", url+"/modules", bytes.NewReader(jsonModule))
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
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       12,
	}
	const secret = "A1B2C3D4E5F6G7H8I9J0K"
	authToken, _ := utils.GenerateToken(user, secret)

	module := fixtures.GetModules()[1]
	const id = "663183f584a7d58b141442ac"
	module.ID, _ = primitive.ObjectIDFromHex(id)

	module.Name = "Updated Module " + time.Now().GoString()
	t.Run("Success Editing ", func(t *testing.T) {

		jsonModule, err := json.Marshal(module)
		if err != nil {
			t.Errorf("unexpected error with %v", err.Error())
		}

		req, errR := http.NewRequest("PUT", url+"/modules", bytes.NewReader(jsonModule))
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
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "Module Updated Successfully", apiResponse.Message, ``)
		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})
	})

	t.Run("Error , Invalid Body Code  400 ", func(t *testing.T) {
		var invalidModule models.Module
		jsonModule, err := json.Marshal(invalidModule)
		if err != nil {
			t.Errorf("unexpected error with %v", err.Error())
		}

		req, errR := http.NewRequest("PUT", url+"/modules", bytes.NewReader(jsonModule))
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
		assert.Equal(t, 400, res.StatusCode, "Bad Request , Can't Edit with invalid body ")
	})
	t.Run("Error , Module Not Found Code 404", func(t *testing.T) {
		unExistObjectID := primitive.NewObjectID()
		module.ID = unExistObjectID
		module.Name = "The New update that does not exist"
		jsonModule, err := json.Marshal(module)
		if err != nil {
			t.Errorf("unexpected error with %v", err.Error())
		}

		req, errR := http.NewRequest("PUT", url+"/modules", bytes.NewReader(jsonModule))
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
		assert.Equal(t, 404, res.StatusCode, "Id Not Found")

	})

}
