package model

type Appauth struct {
	Id       uint   `gorm:"primarykey;not null" json:"id" form:"id"`
	AuthName string `gorm:"not null" json:"auth_name" form:"auth_name"`
	H2Valves int    `gorm:"not null" json:"h2_valves" form:"h2_valves"`
	H2SP     int    `gorm:"not null" json:"h2_sp" form:"h2_sp"`
	AirCmp   int    `gorm:"not null" json:"air_cmp" form:"air_cmp"`
	AirThrot int    `gorm:"not null" json:"air_throt" form:"air_throt"`
	CoolFan  int    `gorm:"not null" json:"cool_fan" form:"cool_fan"`
	CoolET   int    `gorm:"not null" json:"cool_et" form:"cool_et"`
	CoolHT   int    `gorm:"not null" json:"cool_ht" form:"cool_ht"`
}

func (Appauth) TableName() string {
	return "appauth"
}
