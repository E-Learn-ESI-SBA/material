package utils

import (
	"errors"
	"fmt"
	"log"
	"madaurus/dev/material/app/shared"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LightUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	ID       string `json:"id"`
	Avatar   string `json:"avatar"`
	Group    string `json:"group"`
	Year     string `json:"year"`
	Promo    string `json:"promo"`
}
type UserDetails struct {
	LightUser
	jwt.Claims
}

func ParseJwt(signedtoken string, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(signedtoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}

func ValidateToken(signedtoken string, secretKey string) (*UserDetails, error) {
	var user UserDetails

	log.Printf("Getting token .. %v", signedtoken)
	token, err := ParseJwt(signedtoken, secretKey)
	if err != nil {
		log.Printf("Error in parsing token: %v", err.Error())
		return nil, errors.New(shared.INVALID_TOKEN)
	}
	if !token.Valid {
		return nil, errors.New(shared.UNAUTHORIZED)
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {

		log.Printf("Error in converting claims: %v", err.Error())

		return nil, errors.New(shared.INVALID_CREDENTIALS)

	}

	// Case: Token expired
	expTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if expTime.Unix() < time.Now().Local().Unix() {
		return nil, errors.New(shared.EXPIRED_TOKEN)
	}
	user.Email = claims["email"].(string)
	user.Username = claims["username"].(string)
	user.Role = claims["role"].(string)
	user.ID = claims["id"].(string)
	user.Avatar = claims["avatar"].(string)
	user.Group = claims["group"].(string)
	user.Year = claims["year"].(string)
	return &user, nil

}

func GenerateToken(user LightUser, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"id":       user.ID,
		"avatar":   user.Avatar,
		"group":    user.Group,
		"year":     user.Year,
	}
	// add expiration time
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

/*
expTime, _ := claims.GetExpirationTime()

	if expTime.Unix() < time.Now().Local().Unix() {
		return nil, gin.Error{
			Err: errors.New("expired Token"),
		}
	}
*/
