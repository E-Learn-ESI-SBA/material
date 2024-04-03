package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers/courses"
	"madaurus/dev/material/app/handlers/modules"
	"madaurus/dev/material/app/middlewares"
)

func ModuleRoute(g *gin.Engine, collection *mongo.Collection) {
	moduleRoute := g.Group("/modules")
	moduleRoute.GET("/:id", middlewares.Authentication(), modules.GetModuleById(collection))
	moduleRoute.GET("/teacher", middlewares.Authentication(), modules.GetTeacherFilteredModules(collection))
	moduleRoute.GET("/admin", middlewares.Authentication(), courses.GetCoursesByAdmin(collection))
	moduleRoute.POST("/create", middlewares.Authentication(), modules.CreateModule(collection))
	moduleRoute.PUT("/update", middlewares.Authentication(), modules.UpdateModule(collection))
	moduleRoute.DELETE("/delete/:id", middlewares.Authentication(), modules.DeleteModule(collection))
	moduleRoute.PUT("/visibility/:id", middlewares.Authentication(), modules.EditModuleVisibility(collection))
	moduleRoute.POST("/public", middlewares.Authentication(), modules.GetPublicFilteredModules(collection))
	moduleRoute.GET("/public/:id", modules.GetPublicFilteredModules(collection))
}
