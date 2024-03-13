package model

type Paraauth struct {
	Id           uint    `gorm:"primarykey;not null" json:"id" form:"id"`
	ParaName     string  `gorm:"varchar(255);not null" json:"para_name" form:"para_name"`
	MinValue     float64 `gorm:"not null" json:"min_value" form:"min_value"`
	MaxValue     float64 `gorm:"not null" json:"max_value" form:"max_value"`
	ChangeEnable bool    `gorm:"not null" json:"change_enable" form:"change_enable"`
}
