package models

import (
	"github.com/busyLambda/bbkk/domain/server"
	"github.com/busyLambda/bbkk/internal/config"
	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Name    server.ServerName
	Profile config.ServerProfile
}

func NewServer(sp config.ServerProfile, sn server.ServerName, m uint) Server {
	return Server{
		Name:    sn,
		Profile: sp,
	}
}
