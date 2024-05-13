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

	"github.com/stretchr/testify/assert"
)

const url = "http://localhost:8080"

func TestCreateModule(t *testing.T) {
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
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		// Message Assertion
		assert.Equal(t, shared.CREATE_MODULE, apiResponse.Message)

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

		assert.Equal(t, http.StatusNotAcceptable, res.StatusCode)
		assert.Equal(t, shared.INVALID_BODY, apiResponse.Message)
		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})
	})
}

func TestCreateManyModules(t *testing.T) {
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
	authToken, _ := utils.GenerateToken(user, secret)
	modules := fixtures.GetModules()
	t.Run("Success Creating Many Modules", func(t *testing.T) {
		jsonModules, err := json.Marshal(modules)
		if err != nil {
			t.Errorf("UnExpected Error: %v", err)
		}
		req, errR := http.NewRequest("POST", url+"/modules/many", bytes.NewReader(jsonModules))
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
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		// Message Assertion
		assert.Equal(t, shared.CREATE_MODULE, apiResponse.Message)

		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})

	})
	t.Run("Unauthorized Access (Missing Auth Token)", func(t *testing.T) {
		jsonModules, err := json.Marshal(modules)
		if err != nil {
			t.Errorf("UnExpected Error: %v", err)
		}
		req, errR := http.NewRequest("POST", url+"/modules/many", bytes.NewReader(jsonModules))
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
		jsonModules, err := json.Marshal(emptyModule)
		if err != nil {
			t.Errorf("UnExpected Error: %v", err)
		}
		req, errR := http.NewRequest("POST", url+"/modules/many", bytes.NewReader(jsonModules))
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

		assert.Equal(t, http.StatusNotAcceptable, res.StatusCode)
		assert.Equal(t, shared.INVALID_BODY, apiResponse.Message)
		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})

	})
	t.Run("Forbidden Access (Teacher)", func(t *testing.T) {
		user2 := utils.LightUser{
			Email:    "ameri.mohamedayoub@gmail.com",
			Username: "ayoub",
			Role:     "teacher",
			ID:       "12",
			Avatar:   "https://www.google.com",
			Year:     "2021",
			Group:    "1A",
		}
		authToken, _ = utils.GenerateToken(user2, secret)
		jsonModules, err := json.Marshal(modules)
		if err != nil {
			t.Errorf("UnExpected Error: %v", err)
		}
		req, errR := http.NewRequest("POST", url+"/modules/many", bytes.NewReader(jsonModules))
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
		defer t.Cleanup(func() {
			err := res.Body.Close()
			if err != nil {
				return
			}
		})
		assert.Equal(t, http.StatusForbidden, res.StatusCode)
	})

}

func TestGetModuleByTeacher(t *testing.T) {
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
	t.Run("Success Getting Modules By Teacher", func(t *testing.T) {
		req, errR := http.NewRequest("GET", url+"/modules/teacher", nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)

		}
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
	t.Run("Unauthorized Access (Missing Auth Token)", func(t *testing.T) {
		req, errR := http.NewRequest("GET", url+"/modules/teacher", nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)

		}
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	})
	/*	t.Run("Forbidden Access (Teacher)", func(t *testing.T) {
			user2 := utils.LightUser{
				Email:    "qsqs",
				Username: "ayoub",
				Role:     "teacher",
				ID:       "2",
				Avatar:   "https://www.google.com",
			}
			authToken, _ = utils.GenerateToken(user2, secret)
			req, errR := http.NewRequest("GET", url+"/modules/teacher", nil)
			if errR != nil {
				t.Errorf("Error while creating request: %v", errR)

			}
			req.Header.Set("Authorization", "Bearer "+authToken)
			client := &http.Client{}
			res, errS := client.Do(req)
			if errS != nil {
				t.Errorf("Error while getting the response: %v", errS)

			}
			assert.Equal(t, http.StatusForbidden, res.StatusCode)
		})

	*/
}
func TestGetModuleByStudent(t *testing.T) {
	t.Run("Success getting Modules By Student", func(t *testing.T) {

		user := utils.LightUser{
			Email:    "moha@gmail.com",
			Username: "ayoub",
			Role:     "student",
			ID:       "2",
			Avatar:   "https://www.google.com",
			Year:     "2021",
			Group:    "1A",
		}
		const secret = "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"
		authToken, _ := utils.GenerateToken(user, secret)
		req, errR := http.NewRequest("GET", url+"/modules/student", nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)

		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		assert.Equal(t, http.StatusOK, res.StatusCode)

	})

}

// Test The Edit Module
/*
func TestUpdateModule(t *testing.T) {
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "teacher",
		ID:       "12",
	}
	const secret = "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"

	authToken, _ := utils.GenerateToken(user, secret)

	module := fixtures.GetModules()[1]
	const id = "663184f340bb0ad546ce0b30"
	module.ID, _ = primitive.ObjectIDFromHex(id)

	module.Name = "Updated Module " + time.Now().GoString()
	t.Run("Success Editing ", func(t *testing.T) {

		jsonModule, err := json.Marshal(module)
		if err != nil {
			t.Errorf("unexpected error with %v", err.Error())
		}

		req, errR := http.NewRequest("PUT", url+"/modules/"+id, bytes.NewReader(jsonModule))
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

		req, errR := http.NewRequest("PUT", url+"/modules/"+id, bytes.NewReader(jsonModule))
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
		unExistHexId := unExistObjectID.Hex()
		module.ID = unExistObjectID
		module.Name = "The New update that does not exist"
		jsonModule, err := json.Marshal(module)
		if err != nil {
			t.Errorf("unexpected error with %v", err.Error())
		}
		req, errR := http.NewRequest("PUT", url+"/modules/"+unExistHexId, bytes.NewReader(jsonModule))
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
func TestDeleteModule(t *testing.T) {
	// Write Test for the delete module
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       "12",
	}
	user2 := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "teacher",
		ID:       "13",
	}

	const id = "663184f340bb0ad546ce0b30"
	const secret = "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"
	authToken, _ := utils.GenerateToken(user, secret)
	teacherAuthToken, _ := utils.GenerateToken(user2, secret)
	t.Run("Delete :: Success  ", func(t *testing.T) {
		req, errR := http.NewRequest("DELETE", url+"/modules/admin"+id, nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		assert.Equal(t, 200, res.StatusCode, "Module Deleted")
	})
	t.Run("Delete :: Module Note Found ", func(t *testing.T) {
		// ID does not exist , or the user id not the same with teacher_id
		newId := primitive.NewObjectID()
		newHexId := newId.Hex()
		req, errR := http.NewRequest("DELETE", url+"/modules/"+newHexId, nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		assert.Equal(t, 404, res.StatusCode, "Unable to delete the module")
	})
	t.Run("Invalid Id", func(t *testing.T) {
		invalidId := "122"
		req, errR := http.NewRequest("DELETE", url+"/modules/"+invalidId, nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		assert.Equal(t, 400, res.StatusCode, "Unable to delete the module")
	})

	t.Run("Delete :: UNAUTHORIZED", func(t *testing.T) {
		req, errR := http.NewRequest("DELETE", url+"/modules/"+id, nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Authorization", "Bearer "+teacherAuthToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		assert.Equal(t, 401, res.StatusCode, "Unable to delete the module")
	})

}

func TestDeleteModuleByAdmin(t *testing.T) {
	// Write Test for the delete module
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       "12",
	}
	const id = "663184f340bb0ad546ce0b30"
	const secret = "A1B2C3D4E5F6G7H8I9J0K"
	authToken, _ := utils.GenerateToken(user, secret)

	t.Run("Delete :: Success  ", func(t *testing.T) {
		req, errR := http.NewRequest("DELETE", url+"/modules/"+id, nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		assert.Equal(t, 200, res.StatusCode, "Module Deleted")
	})
	t.Run("Delete :: Module Note Found ", func(t *testing.T) {
		// ID does not exist , or the user id not the same with teacher_id
		newId := primitive.NewObjectID()
		newHexId := newId.Hex()
		req, errR := http.NewRequest("DELETE", url+"/modules/"+newHexId, nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		assert.Equal(t, 404, res.StatusCode, "Unable to delete the module")
	})
	t.Run("Invalid Id", func(t *testing.T) {
		id := "122"
		req, errR := http.NewRequest("DELETE", url+"/modules/"+id, nil)
		if errR != nil {
			t.Errorf("Error while creating request: %v", errR)
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("Error while getting the response: %v", errS)
		}
		assert.Equal(t, 400, res.StatusCode, "Unable to delete the module")
	})
}


*/
