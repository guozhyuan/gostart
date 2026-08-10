package router

import (
	"gostart/internal/handler"
	"gostart/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RouteConfig(engine *gin.Engine) {
	api := engine.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/login", handler.Login)
		api.POST("/logout", handler.Logout)
		api.GET("/user", handler.GetUsers)
		api.GET("/user/:id", handler.GetUser)
		api.POST("/regist", handler.Registe)
		api.PUT("/user/:id", handler.UpdateUser)
		api.DELETE("/user/:id", handler.DeleteUser)
	}
}
