package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/enforcement"
	"github.com/permitio/permit-golang/pkg/permit"
	"log"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/shared/iam"
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
		if user.Role == iam.ROLEAdminKey {
			c.Next()
		} else if user.Role == iam.ROLETeacherKey {
			if role == iam.ROLEStudentKey || role == iam.ROLEAdminKey {
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
func StaticRBAC(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := c.Get("user")
		if !err {
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			c.Abort()
			return
		}
		user := claim.(*utils.UserDetails)
		if user.Role != role {
			c.JSON(http.StatusForbidden, gin.H{"message": shared.FORBIDDEN})
			c.Abort()
			return

		}
		c.Next()
		return
	}
}
func IAM(permit *permit.Client, resourceType string, RequestAction string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := c.Get("user")
		if !err {
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			c.Abort()
			return
		}
		moduleId, _ := c.Params.Get("id")

		user := claim.(*utils.UserDetails)
		log.Printf("User Role: %v\n", user.Role)
		module := enforcement.ResourceBuilder(resourceType).WithKey(moduleId).Build()
		EnforceAction := enforcement.Action(RequestAction)
		userRole := enforcement.UserBuilder(user.ID).Build()
		decision, errC := permit.Check(userRole, EnforceAction, module)
		if errC != nil {
			log.Printf("Error While Checking the Permission: %v\n", errC)
			c.JSON(http.StatusForbidden, gin.H{"message": shared.FORBIDDEN})
			c.Abort()
			return
		}
		if !decision {
			c.JSON(http.StatusForbidden, gin.H{"message": shared.FORBIDDEN})
			c.Abort()
			return
		}
		c.Next()
		return

	}

}
