package model

import "gorm.io/gorm"

type Filemap struct {
	gorm.Model
	SoftwareVersionBranch string `gorm:"not null"`
	SoftwareVersionNumber string `gorm:"not null"`
	SoftwareName          string `gorm:"not null"`
	A2lFile               string `gorm:"not null"`
	BuildFile             string `gorm:"not null"`
	OptionalFile          string `gorm:"not null"`
}
