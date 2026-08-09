package model

type Base struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Assemble(code int, message string, data any) *Base {
	return &Base{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
