package common

type LoginResp struct {
	ID           int64  `gorm:"primarykey" json:"id"`
	Username     string `gorm:"size:64;uniqueIndex" json:"username"`
	Email        string `gorm:"size:128" json:"email"`
	Age          int32  `json:"age" gorm:"default:0"`
	AccessToken  string `gorm:"primarykey" json:"accessToken"`
	RefreshToken string `gorm:"primarykey" json:"refreshToken"`
}
type UserResp struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int32  `json:"age"`
}
