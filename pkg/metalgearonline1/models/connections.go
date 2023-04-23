package models

import "gorm.io/gorm"

func init() {
	All = append(All, &Connection{})
}

type Connection struct {
	gorm.Model
	UserID uint
	User   User

	RemoteAddr string `gorm:"type:varchar(15)"`
	RemotePort uint16
	LocalAddr  string `gorm:"type:varchar(15)"`
	LocalPort  uint16
}
