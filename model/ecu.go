package model

type ECUProjectList struct {
	Id             uint   `gorm:"primarykey;not null" json:"id" form:"id"`
	ProjectName    string `gorm:"not null" json:"ProjectName" form:"ProjectName"`
	ProjectCode    string `gorm:"not null" json:"ProjectCode" form:"ProjectCode"`
	SoftwareBranch uint   `gorm:"not null" json:"SoftwareBranch" form:"SoftwareBranch"`
	SoftwareName   string `gorm:"not null" json:"SoftwareName" form:"SoftwareName"`
}

func (ECUProjectList) TableName() string {
	return "ecuprojectlist"
}
