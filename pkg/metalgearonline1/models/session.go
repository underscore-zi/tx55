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
	ID        uuid.UUID `gorm:"type:uuid;primary_key`
	CreatedAt time.Time
	UserID    uint
	User      User
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
