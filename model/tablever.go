package model

type Tablever struct {
	Ver       int    `gorm:"primarykey;not null" json:"Ver" form:"Ver"`
	User      string `gorm:"not null" json:"UserName" form:"UserName"`
	CreatedAt string `gorm:"not null" json:"Time" form:"Time"`
	ChangeLog string `gorm:"not null" json:"ChangeLog" form:"ChangeLog"`
}

func (Tablever) TableName() string {
	return "tablever"
}
