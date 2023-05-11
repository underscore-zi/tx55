package models

import (
	"gorm.io/gorm"
	"time"
)

func init() {
	All = append(All, &Ban{})
}

type BanType byte

const (
	IPBan   BanType = 1
	UserBan BanType = 2
)

func (b BanType) String() string {
	switch b {
	case IPBan:
		return "IP"
	case UserBan:
		return "User"
	default:
		return "Unknown"
	}
}

type Ban struct {
	gorm.Model
	ExpiresAt time.Time
	UserID    uint
	User      User
	CreatedBy string
	UpdatedBy string
	Reason    string
	Type      BanType
}
