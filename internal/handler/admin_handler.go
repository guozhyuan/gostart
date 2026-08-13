package handler

import "github.com/gin-gonic/gin"

// @Summary      管理员登录
// @Description  管理员登录
// @Tags         管理员
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        username  formData string true "用户名"
// @Param        password  formData string true "用户名"
// @Success      200   {object}   model.LoginResp  "登录成功"
// @Failure      400   {object}  model.Base  "请求参数错误"
// @Router       /admin/login [post]
func AdminLogin(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Admin login successful"})
}
