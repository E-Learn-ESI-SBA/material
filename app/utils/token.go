package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserDetails struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Uid      string `json:"uid"`
	jwt.Claims
}

func ValidateToken(signedtoken string, secretKey string) (claims *UserDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &UserDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*UserDetails)
	if !ok {
		msg = "The Token is invalid"
		return
	}
	expTime, _ := claims.GetExpirationTime()
	if expTime.Unix() < time.Now().Local().Unix() {
		msg = "Expired Token"
		return
	}
	return claims, msg
}
