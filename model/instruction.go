package model

type Instruction struct {
	Id                uint   `gorm:"primarykey;not null" json:"id" form:"id"`
	CreateTime        string `gorm:"not null" json:"create_time" form:"create_time"`
	AdminName         string `gorm:"not null" json:"admin_name" form:"admin_name"`
	AdminMechineName  string `gorm:"not null" json:"admin_mechine_name" form:"admin_mechine_name"`
	ClientName        string `gorm:"not null" json:"client_name" form:"client_name"`
	ClientMachineName string `gorm:"not null" json:"client_machine_name" form:"client_machine_name"`
	Result            string `gorm:"not null" json:"result" form:"result"`
}

func (Instruction) TableName() string {
	return "instruction"
}
