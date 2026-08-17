package handler

import (
	model "gostart/internal/dao/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 小对象传值,防止堆分配;大对象传指针,防止内存泄漏
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, model.Assemble(http.StatusOK, "success", data))
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, model.Assemble(code, message, nil))
}
