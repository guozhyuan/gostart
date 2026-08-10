package model

// @Description 统一的 API 响应格式 ,响应成功的结果放在Data字段中;响应失败的结果放在Message字段中
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

type LoginResp struct {
	ID           uint64 `gorm:"primarykey" json:"id"`
	Username     string `gorm:"size:64;uniqueIndex" json:"username"`
	Email        string `gorm:"size:128" json:"email"`
	Age          uint8  `json:"age" gorm:"default:0"`
	AccessToken  string `gorm:"primarykey" json:"accessToken"`
	RefreshToken string `gorm:"primarykey" json:"refreshToken"`
}
type UserResp struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      uint8  `json:"age"`
}
