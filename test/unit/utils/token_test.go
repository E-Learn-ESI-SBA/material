package utils

import (
	"github.com/stretchr/testify/assert"
	"log"
	"madaurus/dev/material/app/utils"
	"testing"
)

func TestValidateToken(t *testing.T) {
	user := utils.LightUser{
		Email:    "ameri.mohamedayoub@gmail.com",
		Username: "ayoub",
		Role:     "admin",
		ID:       "12",
		Avatar:   "sqkddkslqdns",
		Group:    "12",
		Year:     "2021",
	}
	jwt_secret := "aTZ6czFOcTFHekRrZEJHUTB5cFlZZ0M1aXQyR3FiNlltaWx5aDJFUWpIQT0K"
	t.Run("Valid Token", func(t *testing.T) {
		validJwt, err := utils.GenerateToken(user, jwt_secret)
		log.Printf("Token %v", validJwt)
		if err != nil {
			t.Errorf("Error: %v", err.Error())
		}
		claim, err := utils.ValidateToken(validJwt, jwt_secret)
		if err != nil {
			t.Errorf("Error: %v", err.Error())
		}
		assert.NotNil(t, claim)
		assert.IsType(t, &utils.UserDetails{}, claim)
	})
	t.Run("Invalid Token", func(t *testing.T) {
		invalidJwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mNvbSIsImV4cCI6MTcxMjI5NjcyOSwiaWQiOjE2Niwicm9sZSI6ImFkbWluIiwidXNlcm5hbWUiOiJBeW91YiJ9.n-jX7YSdHV1lv_y2Nte-x9CsFxNDUGUI_gVAmh5BQ0k"
		claim, err := utils.ValidateToken(invalidJwt, jwt_secret)

		assert.NotNil(t, err)
		assert.Nil(t, claim)
	})

	t.Run("Short Token", func(t *testing.T) {
		shortJwt := "s1"
		claim, err := utils.ValidateToken(shortJwt, jwt_secret)
		assert.NotNil(t, err)
		assert.Nil(t, claim)
	})
}
