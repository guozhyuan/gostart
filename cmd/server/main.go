package main

import (
	"gostart/internal/config"
	"gostart/internal/pkg"
	"gostart/internal/router"
	"log"
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
	// JWT 密钥必须在启动期就校验：缺失时立即拒绝启动，而不是等到首次登录才暴露。
	// 若此处放行，token 签发/校验将退化为不可用或使用不安全密钥。
	if config.Configs.JWT == nil || config.Configs.JWT.SecretKey == "" {
		log.Fatal("jwt.secretKey 未配置：请在 configs/config.yaml 或环境变量中设置强密钥")
	}
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
