package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/handlers"
	"madaurus/dev/material/app/middlewares"
	"madaurus/dev/material/app/shared/iam"
	"madaurus/dev/material/app/utils"
)

func ModuleRoute(g *gin.Engine, collection *mongo.Collection, permitApi *permit.Client, client *mongo.Client) {
	moduleRoute := g.Group("/modules")
	moduleRoute.GET("/:id", middlewares.Authentication(), handlers.GetModuleById(collection))
	//moduleRoute.GET("/teacher/filter", middlewares.Authentication(), handlers.GetTeacherFilterModules(collection))
	moduleRoute.GET("/teacher", middlewares.Authentication(), handlers.GetModuleByTeacher(collection, permitApi))
	moduleRoute.POST("/", middlewares.Authentication(), middlewares.BasicRBAC("admin"), handlers.CreateModule(collection, client, permitApi))
	moduleRoute.PUT("/:moduleId", middlewares.Authentication(), handlers.UpdateModule(collection))
	// Protected By the admin
	moduleRoute.DELETE("/:id", middlewares.Authentication(), middlewares.BasicRBAC("admin"), handlers.DeleteModule(collection))
	moduleRoute.PUT("/visibility/:id", middlewares.Authentication(), handlers.EditModuleVisibility(collection))
	moduleRoute.GET("/public", middlewares.Authentication(), handlers.GetPublicFilteredModules(collection))
	moduleRoute.GET("/student", middlewares.Authentication(), middlewares.StaticRBAC(iam.ROLEStudentKey), handlers.GetModuleByStudent(collection))
	// Make it public in #Production
	moduleRoute.GET("/public/:id", middlewares.Authentication(), handlers.GetPublicFilteredModules(collection))
	moduleRoute.POST("/many", middlewares.Authentication(), middlewares.BasicRBAC("admin"), handlers.CreateManyModules(collection))
	moduleRoute.POST("/test", middlewares.Authentication(), func(context *gin.Context) {
		context.JSON(200, gin.H{"message": "Hello"})
	})
	moduleRoute.GET("/overview/:id", middlewares.Authentication(), handlers.GetModuleByIdOverview(collection))
	moduleRoute.GET("/admin", middlewares.Authentication(), middlewares.StaticRBAC(iam.ROLEAdminKey), handlers.GetModulesByAdmin(collection))
	moduleRoute.PATCH("/visibility/:id", middlewares.Authentication(), func(context *gin.Context) {
		value, _ := context.Get("user")
		user := value.(*utils.UserDetails)
		utils.GetAllowedResources("delete", "sections", user.ID, permitApi)
	})
}
