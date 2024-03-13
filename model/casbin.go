package model

import "gorm.io/gorm"

type CasinoModel struct {
	gorm.Model
	PType  string `gorm:"column:p_type" json:"p_type" form:"p_type" description:"策略类型"`
	RoleId string `gorm:"column:v0" json:"role_id" form:"v0" description:"角色id"`
	Path   string `gorm:"column:v1" json:"path" form:"v1" description:"api路径"`
	Method string `gorm:"column:v2" json:"method" form:"v2" description:"方法"`
}

func (c *CasinoModel) TableName() string {
	return "casbin_rule"
}

//func (c *CasinoModel) AddPolicy() error {
//	if ok, err := common.CasBin.AddPolicy(c.RoleId, c.Path, c.Method); !ok {
//		return err
//	}
//	return nil
//}
