package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Phone      string `gorm:"primary_key" json:"phone"`
	Avatar     string `json:"avatar"`
	Nickname   string `json:"nickname"`
	CreateTime int64  `json:"create_time"`
}

func (m *User) TableName() string {
	return "user"
}

func initUser(tx *gorm.DB) {
	tx.Create(&User{
		Phone:      "12345678901",
		Avatar:     "https://www.baidu.com",
		Nickname:   "用户1",
		CreateTime: time.Now().Unix(),
	})
}
