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
		api.POST("/login", handler.Login)
		api.POST("/logout", handler.Logout)
		api.GET("/user", handler.GetUsers)
		api.GET("/user/:id", handler.GetUser)
		api.POST("/regist", handler.Registe)
		api.PUT("/user/:id", handler.UpdateUser)
		api.DELETE("/user/:id", handler.DeleteUser)
	}

	admin := engine.Group("/admin")
	admin.Use(middleware.ZapLoggerMiddleware())
	admin.Use(middleware.RateLimitMiddleware(5))
	{
		admin.GET("/login", handler.AdminLogin)
	}

	// Swagger 启用
	if config.Configs.Swagger.Enabled {
		docs.SwaggerInfo.Host = config.Configs.Swagger.Host
		docs.SwaggerInfo.Schemes = config.Configs.Swagger.Scheme
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
