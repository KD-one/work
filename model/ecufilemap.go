package model

type EcuFileMap struct {
	Id          uint   `gorm:"primarykey;not null" json:"id" form:"id"`
	UserName    string `gorm:"not null" json:"user_name" form:"user_name"`
	CreateTime  string `gorm:"not null" json:"create_time" form:"create_time"`
	DeleteTime  string `gorm:"not null" json:"delete_time" form:"delete_time"`
	Branch      uint   `gorm:"not null" json:"branch" form:"branch"`
	Version     uint   `gorm:"not null" json:"version" form:"version"`
	BuildFile   string `gorm:"not null" json:"build_file" form:"build_file"`
	CalMain     uint   `gorm:"not null" json:"cal_main" form:"cal_main"`
	CalSub      uint   `gorm:"not null" json:"cal_sub" form:"cal_sub"`
	CalCmd      string `gorm:"not null" json:"cal_cmd" form:"cal_cmd"`
	ChangeLog   string `gorm:"not null" json:"change_log" form:"change_log"`
	ChangeLogEn string `gorm:"not null" json:"change_log_en" form:"change_log_en"`
}

func (EcuFileMap) TableName() string {
	return "ecufilemap"
}
