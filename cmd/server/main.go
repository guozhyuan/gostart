package main

import (
	"gostart/internal/config"
	"gostart/internal/pkg"
	"gostart/internal/router"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	pkg.ReadConfig()
	pkg.ZapLogInit()
	pkg.ConnectDB()
	pkg.ConnectRedis()

	if config.Configs.Env == "prod" {
		gin.SetMode("release")
	}
	engine := gin.New()
	engine.SetTrustedProxies(config.Configs.Gin.TrustProxy)
	router.RouteConfig(engine)
	s := &http.Server{
		Addr:           config.Configs.Gin.Host,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func Release() {
	// service.ReleaseDB()
	// service.ReleaseRedis()

}
