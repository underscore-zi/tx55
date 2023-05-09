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

type Ban struct {
	ID        uint
	CreatedAt time.Time
	ExpiresAt time.Time
	UserID    uint
	User      User
	Reason    string
	Type      BanType
}
