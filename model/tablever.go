package model

type Tablever struct {
	Ver       int    `gorm:"primarykey;not null" json:"ver" form:"ver"`
	User      string `gorm:"not null" json:"user" form:"user"`
	ChangeLog string `gorm:"not null" json:"change_log" form:"change_log"`
	CreatedAt string `gorm:"not null" json:"created_at" form:"created_at"`
}

func (Tablever) TableName() string {
	return "tablever"
}
