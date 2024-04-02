package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/utils"
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

func GetTeacherFilteredModules(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filterModule interfaces.ModuleFilter
		err := c.BindJSON(&filterModule)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		user := c.MustGet("user").(*utils.UserDetails)
		modules, CursorErr := services.GetModulesByFilter(c.Request.Context(), collection, filterModule, "private", &user.ID)
		if CursorErr != nil {
			c.JSON(400, gin.H{"error": CursorErr.Error()})
			return
		}
		c.JSON(200, gin.H{"modules": modules})
	}
}

func CreateModule(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var module models.Module
		err := c.BindJSON(&module)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		user := c.MustGet("user").(*utils.UserDetails)
		module.TeacherId = user.ID
		err = services.CreateModule(c.Request.Context(), collection, module)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Module Created Successfully"})
	}

}
func UpdateModule(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var module models.Module
		user := c.MustGet("user").(*utils.UserDetails)
		module.TeacherId = user.ID
		err := c.BindJSON(&module)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.UpdateModule(c.Request.Context(), collection, module)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Module Updated Successfully"})
	}
}

func DeleteModule(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		moduleId, errP := c.Params.Get("id")
		if errP != true {
			c.JSON(400, gin.H{"error": errors.New("module ID is Required")})
			return
		}
		err := services.DeleteModule(c.Request.Context(), collection, moduleId, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Module Deleted Successfully"})
	}
}

func GetModuleById(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		moduleId, errP := c.Params.Get("id")
		if errP != true {
			c.JSON(400, gin.H{"error": errors.New("module ID is Required")})
			return
		}
		module, err := services.GetModuleById(c.Request.Context(), collection, moduleId)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"module": module})
	}
}
