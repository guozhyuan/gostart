package handler

import (
	"gostart/internal/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminHandlerInterface interface {
	AdminLogin(c *gin.Context)
}

// @Summary      管理员登录
// @Description  管理员登录
// @Tags         管理员
// @Accept       json
// @Produce      text/plain
// @Success      200   {object}   common.LoginResp  "登录成功"
// @Failure      400   {object}  common.Base  "请求参数错误"
// @Router       /admin/login [get]
func AdminLogin(c *gin.Context) {
	pkg.ZapLogger.Info("Admin login request received", zap.String("hello", "world"))

	c.JSON(200, gin.H{"message": "Admin login successful"})
}
