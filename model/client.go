package model

type ClientLog struct {
	Id        uint   `json:"id" form:"id"`
	UserName  string `json:"UserName" form:"UserName"`
	Time      string `json:"Time" form:"Time"`
	PCName    string `json:"PCName" form:"PCName"`
	ChangeLog string `json:"ChangeLog" form:"ChangeLog"`
}

func (ClientLog) TableName() string {
	return "clientlog"
}
