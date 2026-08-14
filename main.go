package main

import (
	"gostart/internal/config"
	"gostart/internal/router"
	"gostart/internal/service"
	"net/http"
	"time"

	docs "gostart/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gin + Swagger 示例 API
// @version         1.0
// @description     这是一个用 Gin 框架集成的 Swagger 示例
// @host            localhost:8080
// @scheme          https
// @BasePath        /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	service.ReadConfig()
	service.ConnectDB()
	service.ConnectRedis()
	service.ZapLogInit()

	if config.Configs.Env == "prod" {
		gin.SetMode("release")
	}
	engine := gin.New()
	engine.SetTrustedProxies(config.Configs.Gin.TrustProxy)
	router.RouteConfig(engine)
	setupSwag(engine)
	s := &http.Server{
		Addr:           config.Configs.Gin.Host,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func setupSwag(engine *gin.Engine) {
	if config.Configs.Swagger.Enabled {
		docs.SwaggerInfo.Host = config.Configs.Swagger.Host
		docs.SwaggerInfo.Schemes = config.Configs.Swagger.Scheme
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}

func Release() {
	// service.ReleaseDB()
	// service.ReleaseRedis()
}
