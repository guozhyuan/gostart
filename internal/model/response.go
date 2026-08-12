package model

type Base struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 默认用指针,除非结构体极小（≤ 16 字节）
func Assemble(code int, message string, data any) Base {
	return Base{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
