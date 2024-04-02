package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type LightUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	ID       int    `json:"id"`
}
type UserDetails struct {
	LightUser
	Uid string `json:"uid"`
	jwt.Claims
}

func ValidateToken(signedtoken string, secretKey string) (claims *UserDetails, err error) {
	token, err := jwt.ParseWithClaims(signedtoken, &UserDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {

		return nil, gin.Error{Err: err}
	}
	claims, ok := token.Claims.(*UserDetails)
	if !ok {
		return nil, gin.Error{
			Err: errors.New("the Token is invalid"),
		}
	}
	expTime, _ := claims.GetExpirationTime()
	if expTime.Unix() < time.Now().Local().Unix() {
		return nil, gin.Error{
			Err: errors.New("expired Token"),
		}
	}
	return claims, nil
}
