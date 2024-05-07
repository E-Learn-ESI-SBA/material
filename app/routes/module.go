package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
	"madaurus/dev/material/app/utils"
)

func ModuleRoute(g *gin.Engine, collection *mongo.Collection, permit *permit.Client) {
	moduleRoute := g.Group("/modules")
	moduleRoute.GET("/:id", middlewares.Authentication(), handlers.GetModuleById(collection))
	moduleRoute.GET("/teacher/filter", middlewares.Authentication(), handlers.GetTeacherFilterModules(collection))
	moduleRoute.GET("/teacher", middlewares.Authentication(), handlers.GetModuleByTeacher(collection, permit))
	moduleRoute.POST("/", middlewares.Authentication(), middlewares.BasicRBAC("admin"), handlers.CreateModule(collection, permit))
	moduleRoute.PUT("/:moduleId", middlewares.Authentication(), handlers.UpdateModule(collection))
	// Protected By the admin
	moduleRoute.DELETE("/:id", middlewares.Authentication(), middlewares.BasicRBAC("admin"), handlers.DeleteModule(collection))
	moduleRoute.PUT("/visibility/:id", middlewares.Authentication(), handlers.EditModuleVisibility(collection))
	moduleRoute.GET("/public", middlewares.Authentication(), handlers.GetPublicFilteredModules(collection))
	moduleRoute.GET("/public/:id", handlers.GetPublicFilteredModules(collection))
	moduleRoute.POST("/many", middlewares.Authentication(), middlewares.BasicRBAC("admin"), handlers.CreateManyModules(collection))
	moduleRoute.POST("/test/:id", middlewares.Authentication(), middlewares.IAM(permit, "sections", "delete"), func(context *gin.Context) {
		context.JSON(200, gin.H{"message": "Hello"})
	})
	moduleRoute.PATCH("/visibility/:id", middlewares.Authentication(), func(context *gin.Context) {
		value, _ := context.Get("user")
		user := value.(*utils.UserDetails)
		utils.GetAllowedResources("delete", "sections", user.ID, permit)
	})
}
