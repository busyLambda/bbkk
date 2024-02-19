package server

import (
	"fmt"
	"sync"
)

// TODO: Possibly make locks so that we don't have two clients opening em up
type ServerManager struct {
  Servers map[string]*McServer
  Wg *sync.WaitGroup
}

func (sm *ServerManager) GetServer(id string) *McServer {
  return sm.Servers[id]
}

func (sm *ServerManager) AddServer(id string, server *McServer) {
  sm.Servers[id] = server
}

func (sm *ServerManager) StartServer(id string) error {
  s := sm.GetServer(id)
  if s == nil {
    return fmt.Errorf("Server not found")
  }

  go s.Start(sm.Wg)

  s.SetStdout()
  s.SetStdin()

  return nil
}

func (sm *ServerManager) ReadStdout(id string, c chan string) {
  s := sm.GetServer(id)

  if s.Stdout == nil {
    s.SetStdout()
  }

  s.ReadStdout(c)
}

func (sm *ServerManager) WriteStdout(id string, r rune) error {
  s := sm.GetServer(id)

  if s.Stdin == nil {
    s.SetStdin()
  }

  return s.WriteStdin(r)
}
