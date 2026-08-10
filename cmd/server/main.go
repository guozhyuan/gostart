package main

import (
	"gostart/internal/config"
	"gostart/internal/middleware"
	"gostart/internal/router"
	"gostart/internal/service"
	"net/http"
	"time"

	_ "gostart/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gin + Swagger 示例 API
// @version         1.0
// @description     这是一个用 Gin 框架集成的 Swagger 示例
// @host            localhost:8848
// @BasePath        /
func main() {
	service.ReadConfig()
	service.ConnectDB()
	service.ConnectRedis()
	service.ZapLogInit()

	if config.Configs.Gin.Mode == "release" {
		gin.SetMode("release")
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	// engine.Use(gin.Logger())
	engine.Use(middleware.ZapLoggerMiddleware())
	engine.Use(middleware.CorsMiddleware())
	router.RouteConfig(engine)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	s := &http.Server{
		Addr:           config.Configs.Tcp.Host + ":" + config.Configs.Tcp.Port,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}
