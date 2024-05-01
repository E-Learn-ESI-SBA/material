package middlewares

import (
	"github.com/gin-gonic/gin"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
)

func BasicRBAC(role  string )  gin.HandlerFunc {
		return func(c *gin.Context) {

			claim, err := c.Get("user")
			if(err) {
				c.JSON(http.StatusInternalServerError, shared.USER_NOT_INJECTED)

			}
			user := claim.(utils.UserDetails)
			if user.Role != role {
				c.JSON(http.StatusUnauthorized, gin.H{"error": errToken.Error()})

			}
			if errToken != nil {
				c.Abort()
				return
			}
			c.Set("user", claims)
			c.Next()
		}
}