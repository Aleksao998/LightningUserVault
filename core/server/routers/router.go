package routers

import (
	"github.com/Aleksao998/LightingUserVault/core/cache"
	docs "github.com/Aleksao998/LightingUserVault/core/docs"
	userHandler "github.com/Aleksao998/LightingUserVault/core/server/handlers/user"
	"github.com/Aleksao998/LightingUserVault/core/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Config struct {
	CacheEnabled bool
}

// InitRouter initializes a new Gin router with predefined routes and middleware
func InitRouter(logger *zap.Logger, vault storage.Storage, cache cache.Cache, config Config) *gin.Engine {
	r := gin.New()

	// Get global Monitor object
	m := ginmetrics.GetMonitor()

	// Set middleware for gin
	m.Use(r)

	// Set middleware for cors
	r.Use(cors.Default())

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger setup
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	handlerConfig := userHandler.Config{
		CacheEnabled: config.CacheEnabled,
	}

	// Init User Handler
	handler := userHandler.NewUserHandler(logger, vault, cache, handlerConfig)

	// User routes
	userGroup := r.Group("/user")
	{
		userGroup.GET("/:id", handler.GetHandler)
		userGroup.POST("/", handler.SetHandler)
	}

	return r
}
