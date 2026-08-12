package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

/*
	使用 go.uber.org/ratelimit 库来实现限流功能。
	go get go.uber.org/ratelimit@v0.3.1
*/

/*
	使用 github.com/go-redis/redis_rate  未实现
	go get github.com/go-redis/redis_rate/v10
*/

// rate: 每秒允许的请求数 (RPS)
func RateLimitMiddleware(rate int) gin.HandlerFunc {
	// 1. 创建限流器，例如 New(5) 表示每秒允许 5 个请求
	limiter := ratelimit.New(rate)

	return func(c *gin.Context) {
		// 2. 在每次请求处理前，调用 Take() 方法
		// 如果请求速率过快，Take() 会阻塞等待，直到符合速率要求
		now := limiter.Take()

		// 可选：记录等待时间，便于监控
		if waitTime := time.Since(now); waitTime > 0 {
			// 表示当前请求被限流器阻塞了 waitTime 时间
			// 可以在这里记录日志或监控指标
			// log.Printf("Request throttled, waited %v", waitTime)
		}

		// 3. 通过限流检查，继续处理请求
		c.Next()
	}
}
