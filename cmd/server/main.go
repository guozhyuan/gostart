package main

import (
	"myserver/internal/config"
	"myserver/internal/middleware"
	"myserver/internal/router"
	"myserver/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	service.ReadConfig()
	service.ConnectDB()
	service.ConnectRedis()
	service.ZapLogInit()

	// gin.SetMode("release")
	engine := gin.New()
	engine.Use(gin.Recovery())
	// engine.Use(gin.Logger())
	engine.Use(middleware.ZapLoggerMiddleware())
	engine.Use(middleware.CorsMiddleware())
	router.RouteConfig(engine)

	s := &http.Server{
		Addr:           config.Configs.Tcp.Host + ":" + config.Configs.Tcp.Port,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}
