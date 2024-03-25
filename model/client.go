package model

type ClientLog struct {
	Id        uint   `gorm:"primarykey;not null" json:"id" form:"id"`
	UserName  string `gorm:"varchar(255);not null" json:"user_name" form:"user_name"`
	Time      string `gorm:"not null" json:"time" form:"time"`
	PCName    string `gorm:"not null" json:"pc_name" form:"pc_name"`
	ChangeLog string `gorm:"not null" json:"change_log" form:"change_log"`
}

func (ClientLog) TableName() string {
	return "clientlog"
}

// PostClientLogModel 接受请求参数
type PostClientLogModel struct {
	UserName  string `json:"UserName" form:"UserName"`
	Time      string `json:"Time" form:"Time"`
	PCName    string `json:"PCName" form:"PCName"`
	ChangeLog string `json:"ChangeLog" form:"ChangeLog"`
}
