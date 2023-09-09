package routers

import (
	"github.com/Aleksao998/LightingUserVault/core/cache"
	docs "github.com/Aleksao998/LightingUserVault/core/docs"
	userHandler "github.com/Aleksao998/LightingUserVault/core/server/handlers/user"
	"github.com/Aleksao998/LightingUserVault/core/storage"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter initializes a new Gin router with predefined routes and middleware
func InitRouter(vault storage.Storage, cache cache.Cache) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger setup
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Init User Handler
	handler := userHandler.NewUserHandler(vault, cache)

	// User routes
	userGroup := r.Group("/user")
	{
		userGroup.GET("/:id", handler.GetHandler)
		userGroup.POST("/", handler.SetHandler)
	}

	return r
}
