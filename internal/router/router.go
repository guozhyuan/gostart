package router

import (
	"gostart/docs"
	"gostart/internal/config"
	"gostart/internal/handler"
	"gostart/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RouteConfig(engine *gin.Engine) {
	engine.Use(gin.Recovery())
	engine.Use(middleware.CorsMiddleware())

	api := engine.Group("/api")
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.ZapLoggerMiddleware())
	{
		api.POST("/test", handler.UserHandler.Test)
		api.POST("/login", handler.UserHandler.Login)
		api.POST("/logout", handler.UserHandler.Logout)
		api.GET("/user", handler.UserHandler.GetUsers)
		api.GET("/user/:id", handler.UserHandler.GetUser)
		api.POST("/regist", handler.UserHandler.Registe)
		api.PUT("/user/:id", handler.UserHandler.UpdateUser)
		api.DELETE("/user/:id", handler.UserHandler.DeleteUser)
		api.POST("/streamer", handler.UserHandler.GetStreamers)
	}

	admin := engine.Group("/admin")
	admin.Use(middleware.ZapLoggerMiddleware())
	admin.Use(middleware.RateLimitMiddleware(5))
	{
		admin.GET("/login", handler.AdminHandler.AdminLogin)
	}

	// Swagger 启用
	if config.Configs.Swagger.Enabled {
		docs.SwaggerInfo.Host = config.Configs.Swagger.Host
		docs.SwaggerInfo.Schemes = config.Configs.Swagger.Scheme
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
