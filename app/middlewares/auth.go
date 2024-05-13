package middlewares

import (
	"fmt"
	"log"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
				c.JSON(http.StatusUnauthorized, gin.H{"message": shared.UNAUTHORIZED})
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
		log.Printf("Token ", ClientToken)
		c.Set("user", claims)
		c.Next()
	}
}
