package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func init() {
	All = append(All, &Session{})
}

type Session struct {
	ID        uuid.UUID `gorm:"primary_key"`
	CreatedAt time.Time
	UserID    uint
	User      User
}

func (s *Session) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
