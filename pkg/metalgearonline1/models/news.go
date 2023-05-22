package models

import (
	"gorm.io/gorm"
)

func init() {
	All = append(All, &News{})
}

type News struct {
	gorm.Model
	Topic string `gorm:"type:varchar(63)"`
	Body  string `gorm:"type:varchar(512)"`
}
