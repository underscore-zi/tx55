package models

import (
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
	ID        uint
	CreatedAt time.Time
	ExpiresAt time.Time
	UserID    uint
	User      User
	Reason    string
	Type      BanType
}
