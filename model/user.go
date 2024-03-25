package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"varchar(20);not null" json:"name" form:"name"`
	Password  string `gorm:"size:255;not null" json:"password" form:"password"`
	UserLevel int    `gorm:"not null" json:"user_level" form:"user_level"`
	AppAuth   string `gorm:"varchar(255);not null" json:"app_auth" form:"app_auth"`
	ParaAuth  string `gorm:"varchar(255);not null" json:"para_auth" form:"para_auth"`
}

type ClientList struct {
	InstructionId     uint   `json:"InstructionId"`
	InstructionResult string `json:"InstructionResult"`
}
