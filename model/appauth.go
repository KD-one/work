package model

type Appauth struct {
	Id       uint   `gorm:"primarykey;not null" json:"id" form:"id"`
	AuthName string `gorm:"not null" json:"AuthName" form:"AuthName"`
	H2Valves int    `gorm:"not null" json:"H2Valves" form:"H2Valves"`
	H2SP     int    `gorm:"not null" json:"H2SP" form:"H2SP"`
	AirCmp   int    `gorm:"not null" json:"AirCmp" form:"AirCmp"`
	AirThrot int    `gorm:"not null" json:"AirThrot" form:"AirThrot"`
	CoolFan  int    `gorm:"not null" json:"CoolFan" form:"CoolFan"`
	CoolET   int    `gorm:"not null" json:"CoolET" form:"CoolET"`
	CoolHT   int    `gorm:"not null" json:"CoolHT" form:"CoolHT"`
}

func (Appauth) TableName() string {
	return "appauth"
}
