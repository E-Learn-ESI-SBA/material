package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
)

func ModuleRoute(g *gin.Engine, collection *mongo.Collection, permit *permit.Client) {
	moduleRoute := g.Group("/modules")
	moduleRoute.GET("/:id", middlewares.Authentication(), handlers.GetModuleById(collection))
	moduleRoute.GET("/teacher", middlewares.Authentication(), handlers.GetTeacherFilterModules(collection))
	moduleRoute.POST("/", middlewares.Authentication(), middlewares.BasicRBAC("admin"), handlers.CreateModule(collection))
	moduleRoute.PUT("/:moduleId", middlewares.Authentication(), handlers.UpdateModule(collection))
	// Protected By the admin
	moduleRoute.DELETE("/:id", middlewares.Authentication(), middlewares.BasicRBAC("admin"), handlers.DeleteModule(collection))
	moduleRoute.PUT("/visibility/:id", middlewares.Authentication(), handlers.EditModuleVisibility(collection))
	moduleRoute.GET("/public", middlewares.Authentication(), handlers.GetPublicFilteredModules(collection))
	moduleRoute.GET("/public/:id", handlers.GetPublicFilteredModules(collection))
	moduleRoute.POST("/many", middlewares.Authentication(), middlewares.BasicRBAC("admin"), handlers.CreateManyModules(collection))
	moduleRoute.POST("/test/:id", middlewares.Authentication(), middlewares.IAM(permit, "chapters", "delete"), func(context *gin.Context) {
		context.JSON(200, gin.H{"message": "Hello"})
	})
}
