package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID        string `gorm:"type:uuid;primary_key;"`
	UserAgent string
	UserID    uint
}

func NewSession(userID uint, userAgent string) Session {
	return Session{
		ID:        uuid.NewString(),
		UserID:    userID,
		UserAgent: userAgent,
	}
}
