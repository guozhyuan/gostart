package model

import "time"

type UserDO struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"size:64;uniqueIndex" json:"username"`
	Password  string    `gorm:"size:128" json:"password"`
	Email     string    `gorm:"size:128" json:"email"`
	Age       uint8     `json:"age" gorm:"default:0"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (UserDO) TableName() string {
	return "user"
}

type TokenDO struct {
	AccessToken  string `gorm:"primarykey" json:"accessToken"`
	RefreshToken string `gorm:"primarykey" json:"refreshToken"`
}
