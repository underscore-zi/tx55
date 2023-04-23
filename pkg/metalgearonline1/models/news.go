package models

import (
	"gorm.io/gorm"
	"time"
)

func init() {
	All = append(All, &News{})
}

type News struct {
	gorm.Model
	Time  time.Time
	Topic string `gorm:"type:varchar(63)"`
	Body  string `gorm:"type:varchar(899)"`
}
