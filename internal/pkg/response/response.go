package response

type LoginResp struct {
	ID           uint64 `gorm:"primarykey" json:"id"`
	Username     string `gorm:"size:64;uniqueIndex" json:"username"`
	Email        string `gorm:"size:128" json:"email"`
	Age          uint8  `json:"age" gorm:"default:0"`
	AccessToken  string `gorm:"primarykey" json:"accessToken"`
	RefreshToken string `gorm:"primarykey" json:"refreshToken"`
}
type UserResp struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      uint8  `json:"age"`
}
