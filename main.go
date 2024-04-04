package main

import (
	"fmt"
	"log"
	"madaurus/dev/material/app/utils"
)

func main() {
	jwt_secret := "A1B2C3D4E5F6G7H8I9J0K"
	validJwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImF5b3ViQGdtYWlsLmNvbSIsImV4cCI6MTcxMjI5NjcyOSwiaWQiOjE2Niwicm9sZSI6ImFkbWluIiwidXNlcm5hbWUiOiJBeW91YiJ9.n-jX7YSdHV1lv_y2Nte-x9CsFxNDUGUI_gVAmh5BQ0k"
	token, err := utils.ValidateToken(validJwt, jwt_secret)
	if err != nil {
		println("HERE Error .........!")
		log.Fatal(err)
	}
	fmt.Println(token.Email)
}
