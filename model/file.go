package model

import "time"

type File struct {
	Uuid          int       `gorm:"primarykey"`
	Name          string    `gorm:"not null"`
	Type          string    `gorm:"not null"`
	Size          string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	ProjectNumber string    `gorm:"not null"`
	VersionNumber string    `gorm:"not null"`
}
