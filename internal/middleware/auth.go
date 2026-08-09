package middleware

import (
	"myserver/internal/handler"
	"myserver/internal/pkg"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

var whitelist = []string{
	"/api/login",
	"/api/regist",
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if slices.Contains(whitelist, path) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.Fail(c, http.StatusUnauthorized, "缺少 Authorization Header")
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			handler.Fail(c, http.StatusUnauthorized, "Authorization Header 格式不正确")
			c.Abort()
			return
		}

		claims, err := pkg.ParseAndValidateToken(parts[1], "access")
		if err != nil {
			handler.Fail(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("userId", claims.UserId)
		c.Next()
	}
}
