package handlers

import (
	"errors"
	"log"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Edit Module Visibility
// @Description Protected Route used to edit module visibility
// @Produce json
// @Accept json
// @Tags Modules
// @Param id path string true "Module ID"
// @Param visibility query string true "Module Visibility"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /modules/visibility/{id} [PUT]
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

// @Summary Get Public Modules
// @Description Protected Route used to get public modules
// @Produce json
// @Accept json
// @Tags Modules
// @Param filter body interfaces.ModuleFilter true "Module Filter"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /modules/public [POST]
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

// @Summary Get Teacher Modules
// @Description Protected Route used to get teacher modules
// @Produce json
// @Tags Modules
// @Accept json
// @Param filter body interfaces.ModuleFilter true "Module Filter"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /modules/teacher [GET]
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

// @Summary Create Module
// @Description Protected Route used to create a module
// @Produce json
// @Accept json
// @Tags Modules
// @Param module body models.Module true "Module Object"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /modules/create [POST]
func CreateModule(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		print("Create Module Handler ...")
		var module models.Module

		err := c.ShouldBindJSON(&module)
		if err != nil {
			log.Println(err.Error())
			c.JSON(400, gin.H{"error": errors.New("invalid Module Object").Error()})
			return
		}
		user := c.MustGet("user").(*utils.UserDetails)
		module.TeacherId = user.ID
		err = services.CreateModule(c.Request.Context(),collection, module)
		if err != nil {
			log.Println(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Module Created Successfully"})
	}

}

// @Summary Update Module
// @Description Protected Route used to update a module
// @Produce json
// @Accept json
// @Tags Modules
// @Param module body models.Module true "Module Object"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /modules/update [PUT]
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

// @Summary Delete Module
// @Description Protected Route used to delete a module
// @Produce json
// @Accept json
// @Param id path string true "Module ID"
// @Success 200 {object} interfaces.APiSuccess
// @Tags Modules
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /modules/delete/{id} [DELETE]
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

// @Summary Get Module By ID
// @Description Get Module By ID
// @Produce json
// @Accept json
// @Tags Modules
// @Param id path string true "Module ID"
// @Success 200 {object} models.ExtendedModule
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /modules/{id} [GET]
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
