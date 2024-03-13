package model

type Appauth struct {
	Id       uint   `gorm:"primarykey;not null" json:"id" form:"id"`
	Name     string `gorm:"not null" json:"name" form:"name"`
	Login    bool   `gorm:"not null" json:"login" form:"login"`
	Register bool   `gorm:"not null" json:"register" form:"register"`
	Update   bool   `gorm:"not null" json:"update" form:"update"`
}

func (Appauth) TableName() string {
	return "appauth"
}
