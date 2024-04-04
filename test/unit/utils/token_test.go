package utils

import (
	"github.com/stretchr/testify/assert"
	"madaurus/dev/material/app/utils"
	"testing"
)

func TestValidateToken(t *testing.T) {
	//jwt_secret := os.Getenv("JWT_SECRET")
	jwt_secret := "A1B2C3D4E5F6G7H8I9J0K"
	t.Run("Valid Token", func(t *testing.T) {
		validJwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImF5b3ViQGdtYWlsLmNvbSIsImV4cCI6MTcxMjI5NjcyOSwiaWQiOjE2Niwicm9sZSI6ImFkbWluIiwidXNlcm5hbWUiOiJBeW91YiJ9.n-jX7YSdHV1lv_y2Nte-x9CsFxNDUGUI_gVAmh5BQ0k"
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

}
