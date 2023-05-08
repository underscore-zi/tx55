package models

import "gorm.io/gorm"

func init() {
	All = append(All, &Notification{})
}

type Notification struct {
	gorm.Model
	UserID uint
	User   User
	// IsImportant comes from a leaked MGS4 elf with debug symbols, but it doesn't seem to have an impact
	IsImportant bool
	HasRead     bool
	Title       string `gorm:"type:varchar(64)"`
	Body        string `gorm:"type:varchar(512)"`
}
