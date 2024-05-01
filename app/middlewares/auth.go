package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"madaurus/dev/material/app/utils"
	"net/http"
	"os"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := os.Getenv("JWT_SECRET")
		ClientToken := c.Request.Header.Get("Authorization")

		var err error
		if len(ClientToken) < 16 {
			// try to get it from cookies
			ClientToken, err = c.Cookie("accessToken")
			if err != nil || ClientToken == "" {
				fmt.Println("No Authorization Header Provided")
				c.JSON(http.StatusBadRequest, gin.H{"message": "No Authorization Header Provided"})
				c.Abort()
				return
			}
		}
		// remove Bearar from token
		ClientToken = ClientToken[7:]

		claims, errToken := utils.ValidateToken(ClientToken, secretKey)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": errToken.Error()})
			c.Abort()
			return
		}
		c.Set("user", claims)
		c.Next()
	}
}
