package middlewares

import (
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
		if ClientToken == "" {
			// try to get it from cookies
			ClientToken, err = c.Cookie("accessToken")
			if err != nil || ClientToken == "" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
				c.Abort()
				return
			}
		}
		claims, errToken := utils.ValidateToken(ClientToken, secretKey)
		if errToken != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("user", claims)
		c.Next()
	}
}
