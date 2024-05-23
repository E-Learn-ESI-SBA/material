package handlers

import (
	"errors"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
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
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /modules/public [POST]
func GetPublicFilteredModules(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		modules, CursorErr := services.GetModulesByFilter(c.Request.Context(), collection)
		if CursorErr != nil {
			log.Println("Error While Getting Public Modules %v", CursorErr)
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.UNABLE_GET_MODULES})
			return

		}
		c.JSON(http.StatusOK, gin.H{"data": modules})
	}
}

/*
// @Summary Get Teacher Modules
// @Description Protected Route used to get teacher modules
// @Produce json
// @Tags Modules
// @Accept json
// @Param filter body interfaces.ModuleFilter true "Module Filter"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /modules/teacher [GET]
func GetTeacherFilterModules(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filterModule interfaces.ModuleFilter
		err := c.ShouldBind(&filterModule)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		value, notFound := c.Get("user")
		if notFound != true {
			c.JSON(400, gin.H{"error": errors.New("user not found")})
			return
		}
		user := value.(*utils.UserDetails)
		modules, CursorErr := services.GetModulesByFilter(c.Request.Context(), collection, filterModule, "private", &user.ID)
		if CursorErr != nil {
			c.JSON(400, gin.H{"error": CursorErr.Error()})
			return
		}
		c.JSON(200, gin.H{"modules": modules})
	}
}
*/
// @Summary Get Module By Teacher
// @Description Protected Route used to get modules by teacher
// @Produce json
// @Accept json
// @Tags Modules
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /modules/teacher [GET]
func GetModuleByTeacher(collection *mongo.Collection, permit *permit.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, notFound := c.Get("user")
		if notFound != true {
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		user := value.(*utils.UserDetails)
		// Change it for #production
		RKeys := utils.GetAllowedResources("read", "modules", user.ID, permit)
		if len(RKeys) == 0 {
			log.Println("The Key is null")
			modules, err := services.GetModulesByTeacher(c.Request.Context(), collection, user.ID)
			if err != nil {
				log.Println("Error While Getting the keys modules %v ", err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"message": shared.UNABLE_GET_MODULE})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": modules})
		} else {
			log.Printf("Keys: %v\n", RKeys)
			log.Println("Here")
			modules, err := services.ModuleSelector(c.Request.Context(), collection, RKeys)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": shared.UNABLE_GET_MODULE})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": modules})

		}
	}
}

// @Summary Create Module
// @Description Protected Route used to create a module
// @Produce json
// @Accept json
// @Tags Modules
// @Param module body models.Module true "Module Object"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /modules [POST]
func CreateModule(collection *mongo.Collection, client *mongo.Client, permit *permit.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		print("Create Module Handler ...")
		var module models.Module
		err := c.ShouldBindJSON(&module)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_BODY})
			return
		}
		err = services.CreateModule(c.Request.Context(), collection, module, permit, client)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.UNABLE_CREATE_MODULE})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": shared.CREATE_MODULE})
	}

}

// @Summary Update Module
// @Description Protected Route used to update a module
// @Produce json
// @Accept json
// @Tags Modules
// @Param module body models.Module true "Module Object"
// @Param moduleId path string true "Module Id"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /modules/{moduleId} [PUT]
func UpdateModule(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var module models.Module
		value, notFound := c.Get("user")
		if notFound != true {
			c.JSON(401, gin.H{"error": errors.New("user not found").Error()})
			return
		}
		moduleId, _ := c.Params.Get("moduleId")
		user := value.(*utils.UserDetails)
		module.TeacherId = user.ID
		err := c.ShouldBindJSON(&module)
		module.ID, _ = primitive.ObjectIDFromHex(moduleId)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		err = services.UpdateModule(c.Request.Context(), collection, module)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Module Updated Successfully"})
	}
}

// @Summary Get Module By ID
// @Description Get Module By ID
// @Produce json
// @Accept json
// @Tags Modules
// @Param id path string true "Module ID"
// @Success 200 {object} models.ExtendedModule
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /modules/{id} [GET]
func GetModuleById(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		moduleId, errP := c.Params.Get("id")
		if errP != true {
			c.JSON(400, gin.H{"error": shared.REQUIRED_ID})
			return
		}
		ModuleIObjectId, errD := primitive.ObjectIDFromHex(moduleId)
		if errD != nil {
			c.JSON(400, gin.H{"error": shared.REQUIRED_ID})
			return
		}
		module, err := services.GetModuleById(c.Request.Context(), collection, ModuleIObjectId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": module})
	}
}

// @Summary Delete Module
// @Description Protected Route used to delete a module
// @Produce json
// @Success 200 {object} interfaces.APiSuccess
// @Tags Modules
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /transaction/module/{id} [DELETE]
func DeleteModule(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		moduleId, errP := c.Params.Get("id")

		if errP != true {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		moduleObjectId, errD := primitive.ObjectIDFromHex(moduleId)
		if errD != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_ID})
			return
		}

		err := services.DeleteModule(c.Request.Context(), collection, moduleObjectId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": shared.DELETE_MODULE})
	}
}

// @Summary Create Module
// @Description Protected Route used to create a module
// @Produce json
// @Accept json
// @Tags Modules
// @Param module body []models.Module true "Module Object"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /modules/many [POST]
func CreateManyModules(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var modules []models.Module
		err := c.ShouldBindJSON(&modules)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_BODY})
			return
		}
		err = services.CreateManyModules(c.Request.Context(), collection, modules)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": shared.CREATE_MODULE})
	}
}

// @Summary Get Module By Student
// @Description Protected Route used to get modules by student
// @Produce json
// @param Authorization header string true "Bearer
// @Tags Modules
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /modules/student [GET]
func GetModuleByStudent(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		value, errV := context.Get("user")
		if errV != true {
			context.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		user := value.(*utils.UserDetails)

		modules, err := services.GetModuleByStudent(context.Request.Context(), collection, user.Year)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.UNABLE_GET_MODULE})
			return
		}
		context.JSON(http.StatusOK, gin.H{"data": modules})
	}
}
func GetModuleByIdOverview(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		moduleId, errD := primitive.ObjectIDFromHex(id)
		if errD != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		module, err := services.GetModuleByIdOverview(c.Request.Context(), collection, moduleId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": module})
	}
}
func GetModulesByAdmin(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		modules, err := services.GetModulesByAdmin(context.Request.Context(), collection)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.UNABLE_GET_MODULE})
			return
		}
		context.JSON(http.StatusOK, gin.H{"data": modules})
	}
}
