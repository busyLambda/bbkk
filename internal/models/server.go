package models

import (
	"github.com/busyLambda/bbkk/domain/server"
	"github.com/busyLambda/bbkk/internal/util"
	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Name         server.ServerName
	DedicatedRam uint
}

// TODO: Rework the profile system.
func NewServer(sf *util.ServerForm) Server {
	return Server{
		Name:         sf.Name,
		DedicatedRam: sf.DedicatedRam,
	}
}
