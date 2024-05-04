package middlewares

import (
	"github.com/gin-gonic/gin"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
)

// This Used only for statical actions
func BasicRBAC(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := c.Get("user")
		if !err {
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			c.Abort()
			return
		}
		user := claim.(*utils.UserDetails)
		if user.Role == shared.ADMIN {

			c.Next()
		} else if user.Role == shared.TEACHER {
			if role == shared.STUDENT || role == shared.ADMIN {
				c.JSON(http.StatusForbidden, gin.H{"message": shared.FORBIDDEN})
				c.Abort()
				return
			}
			c.Next()
		}
		c.Next()
		return
	}
}
