package middleware

import (
	"fmt"
	"io"
	"log/slog"
	"myserver/internal/service"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Deprecated: Use LogMiddleware instead.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
		)
	})
}

// zap 输出结构化日志
func ZapLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		proto := c.Request.Proto
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		var msg string
		if status >= 500 {
			msg = "服务器错误"
		} else {
			msg = "请求处理完成"

		}
		service.ZapLogger.Info(msg,
			zap.String("proto", proto),
			zap.String("method", method),
			zap.Int("status", status),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}

// Gin日志中间件
func LogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s\" %s %s\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC3339),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		// SkipPaths: []string{"/api/getStudents"}, // 日志跳过特定路径
		Output: outputWriter(),
	})
}

func outputWriter() io.Writer {
	// os.O_CREATE — 文件不存在时创建;  os.O_APPEND — 文件存在时追加写入;	os.O_WRONLY — 只写模式
	f, _ := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return io.MultiWriter(f, os.Stdout)
}
