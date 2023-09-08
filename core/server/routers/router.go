package routers

import (
	docs "github.com/Aleksao998/LightingUserVault/core/docs"
	"github.com/Aleksao998/LightingUserVault/core/server/handlers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(handler *handlers.UserHandler) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/user/:id", handler.GetUserHandler)
	r.POST("/user", handler.SetUserHandler)

	return r
}
