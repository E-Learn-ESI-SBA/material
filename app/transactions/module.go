package transactions

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
)

// @Summary Delete Module
// @Description Protected Route used to delete a module By admin
// @Produce json
// @Success 200 {object} interfaces.APiSuccess
// @Tags Modules
// @Failure 400 {object} interfaces.APiError
// @Failure 403 {object} interfaces.APiError Not Allowed
// @Failure 500 {object} interfaces.APiError
// @Router /transaction/admin/{id} [DELETE]
func DeleteModuleByAdmin(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		moduleId, errP := c.Params.Get("id")
		if errP != true {
			c.JSON(400, gin.H{"error": errors.New(shared.REQUIRED_ID)})
			return
		}
		id, errId := primitive.ObjectIDFromHex(moduleId)
		if errId != nil {
			c.JSON(400, errors.New(shared.REQUIRED_ID))
			return
		}
		module, err := services.GetModuleById(c.Request.Context(), collection, id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"module": module})
	}
}

// @Summary Delete Module
// @Description Protected Route used to delete a module
// @Produce json
// @Success 200 {object} interfaces.APiSuccess
// @Tags Modules
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /transaction/module/{id} [DELETE]
func DeleteModule(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, notFound := c.Get("user")
		if notFound {
			c.JSON(401, gin.H{"error": errors.New(shared.USER_NOT_INJECTED)})
			return
		}
		user := value.(*utils.UserDetails)
		moduleId, errP := c.Params.Get("id")

		if errP != true {
			c.JSON(400, gin.H{"error": errors.New(shared.REQUIRED_ID)})
			return
		}
		err := services.DeleteModule(c.Request.Context(), collection, moduleId, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": shared.DELETE_MODULE})
	}
}
