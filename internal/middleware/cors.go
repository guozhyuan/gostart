package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	// cors.Default()
	handlerfunc := cors.New(cors.Config{
		AllowOrigins:     []string{"http://172.31.98.14", "http://172.31.98.14", "https://localhost:9090", "http://localhost:9090"}, // 明确允许 HTTPS 来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	return handlerfunc
}
