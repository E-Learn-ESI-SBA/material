package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type LightUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	ID       string `json:"id"`
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
	token, err := ParseJwt(signedtoken, secretKey)
	if err != nil {
		return nil, errors.New("invalid Token")
	}
	if !token.Valid {
		return nil, errors.New("UNAUTHORIZED")
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("invalid Token")

	}

	// Case: Token expired
	expTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if expTime.Unix() < time.Now().Local().Unix() {
		return nil, errors.New("expired Token")
	}
	user.Email = claims["email"].(string)
	user.Username = claims["username"].(string)
	user.Role = claims["role"].(string)
	user.ID = claims["id"].(string)
	return &user, nil

}

func GenerateToken(user LightUser, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"id":       user.ID,
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
func TestToken() (error, string) {
	// Generate new token
	token, err := GenerateToken(LightUser{
		ID:       "12",
		Email:    "ameri.ayoub@gmail.com",
		Username: "Ayoub",
		Role:     "admin",
	}, "A1B2C3D4E5F6G7H8I9J0K")

	return err, token
}
