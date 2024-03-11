package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	UserAgent string
	UserID    uint
}

func NewSession(userID uint, userAgent string) Session {
	return Session{
		UserID:    userID,
		UserAgent: userAgent,
	}
}
