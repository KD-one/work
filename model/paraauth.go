package model

type Paraauth struct {
	Id           uint    `gorm:"primarykey;not null" json:"id" form:"id"`
	ParaName     string  `gorm:"varchar(255);not null" json:"ParaName" form:"ParaName"`
	MinValue     float64 `gorm:"not null" json:"MinValue" form:"MinValue"`
	MaxValue     float64 `gorm:"not null" json:"MaxValue" form:"MaxValue"`
	ChangeEnable bool    `gorm:"not null" json:"ChangeEnable" form:"ChangeEnable"`
}
