package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	handlers "madaurus/dev/material/app/handlers/modules"
	"madaurus/dev/material/app/middlewares"
)

func ModuleRoute(g *gin.Engine, collection *mongo.Collection) {
	moduleRoute := g.Group("/modules")
	moduleRoute.GET("/:id", middlewares.Authentication(), handlers.GetModuleById(collection))
	moduleRoute.GET("/teacher", middlewares.Authentication(), handlers.GetTeacherFilteredModules(collection))
	moduleRoute.POST("/", middlewares.Authentication(), handlers.CreateModule(collection))
	moduleRoute.PUT("/:moduleId", middlewares.Authentication(), handlers.UpdateModule(collection))
	moduleRoute.DELETE("/:id", middlewares.Authentication(), handlers.DeleteModule(collection))
	moduleRoute.PUT("/visibility/:id", middlewares.Authentication(), handlers.EditModuleVisibility(collection))
	moduleRoute.GET("/public", middlewares.Authentication(), handlers.GetPublicFilteredModules(collection))
	moduleRoute.GET("/public/:id", handlers.GetPublicFilteredModules(collection))
}
