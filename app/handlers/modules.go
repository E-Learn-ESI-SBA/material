package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/services"
)

func EditModuleVisibility(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// module should be body type {id: string, isPublic: bool}
		id, errP := c.Params.Get("id")
		visibility := c.Query("visibility")
		if errP != true || visibility == "" {
			c.JSON(400, gin.H{"error": errors.New("module ID is Required")})
			return
		}
		err := services.EditModuleVisibility(c.Request.Context(), collection, id, visibility == "visible")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Module Visibility Updated Successfully"})
	}
}

func GetPublicFilteredModules(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filterModule interfaces.ModuleFilter
		err := c.BindJSON(&filterModule)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		modules, CursorErr := services.GetModulesByFilter(c.Request.Context(), collection, filterModule, "public", nil)
		if CursorErr != nil {
			c.JSON(400, gin.H{"error": CursorErr.Error()})
			return

		}
		c.JSON(200, gin.H{"modules": modules})
	}
}
