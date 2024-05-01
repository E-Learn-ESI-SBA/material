package integrating

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"madaurus/dev/material/app/interfaces"
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
	t.Run("Success Creating Module", func(t *testing.T) {

		module := fixtures.GetModules()[0]
		jsonModule, err := json.Marshal(module)
		if err != nil {
			t.Errorf("UnExpected Erro: " + err.Error())
		}
		req, errR := http.NewRequest("POST", url+"module", bytes.NewReader(jsonModule))
		if errR != nil {
			t.Errorf("UnExpected Erro: " + errR.Error())

		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		client := &http.Client{}
		res, errS := client.Do(req)
		if errS != nil {
			t.Errorf("UnExpected Erro: " + errS.Error())

		}
		body, err := io.ReadAll(res.Body)
		t.Log("This is The Body", body)
		response := interfaces.APiSuccess{Message: "Module Deleted Successfully", Code: 200}
		mockJson, _ := json.Marshal(response)
		assert.Equal(t, mockJson, body)
	})
}
