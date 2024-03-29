package models

import (
	"github.com/busyLambda/bbkk/domain/server"
	"github.com/busyLambda/bbkk/internal/config"
	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Name server.ServerName
}

// TODO: Rework the profile system.
func NewServer(sp config.ServerProfile, sn server.ServerName, m uint) Server {
	return Server{
		Name: sn,
	}
}
