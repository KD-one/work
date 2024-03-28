package model

type ECUVer struct {
	Id                 uint   `gorm:"primarykey;not null" json:"id" form:"id"`
	ChangeInitiator    string `gorm:"not null" json:"ChangeInitiator" form:"ChangeInitiator"`
	ChangeInitiateTime string `gorm:"not null" json:"ChangeInitiateTime" form:"ChangeInitiateTime"`
	ChangeCause        string `gorm:"not null" json:"ChangeCause" form:"ChangeCause"`
	ChangeReq          string `gorm:"not null" json:"ChangeReq" form:"ChangeReq"`
	ChangeAttached     string `gorm:"not null" json:"ChangeAttached" form:"ChangeAttached"`
	ChangeApplyRange   string `gorm:"not null" json:"ChangeApplyRange" form:"ChangeApplyRange"`
	SWModifier         string `gorm:"not null" json:"SWModifier" form:"SWModifier"`
	SWFinishTime       string `gorm:"not null" json:"SWFinishTime" form:"SWFinishTime"`
	SWLog              string `gorm:"not null" json:"SWLog" form:"SWLog"`
	SWBuildFile        string `gorm:"not null" json:"SWBuildFile" form:"SWBuildFile"`
	SWBranch           uint   `gorm:"not null" json:"SWBranch" form:"SWBranch"`
	SWVersion          uint   `gorm:"not null" json:"SWVersion" form:"SWVersion"`
	SWCalMain          uint   `gorm:"not null" json:"SWCalMain" form:"SWCalMain"`
	SWCalSub           uint   `gorm:"not null" json:"SWCalSub" form:"SWCalSub"`
	HILTester          string `gorm:"not null" json:"HILTester" form:"HILTester"`
	HILFinishTime      string `gorm:"not null" json:"HILFinishTime" form:"HILFinishTime"`
	HILResult          string `gorm:"not null" json:"HILResult" form:"HILResult"`
	SysVerifier        string `gorm:"not null" json:"SysVerifier" form:"SysVerifier"`
	SysVerifyTime      string `gorm:"not null" json:"SysVerifyTime" form:"SysVerifyTime"`
	SysVerifyResult    string `gorm:"not null" json:"SysVerifyResult" form:"SysVerifyResult"`
	SysVerifyAttached  string `gorm:"not null" json:"SysVerifyAttached" form:"SysVerifyAttached"`
	SWLevel            uint   `gorm:"not null" json:"SWLevel" form:"SWLevel"`
	SWLogClient        string `gorm:"not null" json:"SWLogClient" form:"SWLogClient"`
	SWLogClientEN      string `gorm:"not null" json:"SWLogClientEN" form:"SWLogClientEN"`
	CALFile            string `gorm:"not null" json:"CALFile" form:"CALFile"`
	Comment            string `gorm:"not null" json:"Comment" form:"Comment"`
}

func (ECUVer) TableName() string {
	return "ecusoftwareversion"
} //ljp-限定表名称只能是这个，不能有加s等任何修改
