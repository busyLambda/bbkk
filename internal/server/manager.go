package server

import (
	"fmt"
	"sync"
)

// TODO: Possibly make locks so that we don't have two clients opening em up
type ServerManager struct {
	Servers map[uint]*McServer
	Wg      *sync.WaitGroup
}

func NewServerManager() *ServerManager {
	return &ServerManager{
		Servers: make(map[uint]*McServer),
		Wg:      &sync.WaitGroup{},
	}
}

func (sm *ServerManager) GetServer(id uint) *McServer {
	return sm.Servers[id]
}

func (sm *ServerManager) AddServer(id uint, server *McServer) {
	sm.Servers[id] = server
}

func (sm *ServerManager) StartServer(id uint) error {
	s := sm.GetServer(id)
	if s == nil {
		return fmt.Errorf("server not found")
	}

	go s.Start(sm.Wg)

	s.SetStdout()
	s.SetStdin()

	return nil
}

func (sm *ServerManager) ReadStdout(id uint, c chan string) {
	s := sm.GetServer(id)

	s.ReadStdout(c)
}

func (sm *ServerManager) WriteStdout(id uint, r rune) error {
	s := sm.GetServer(id)

	if s.Stdin == nil {
		s.SetStdin()
	}

	return s.WriteRune(r)
}
