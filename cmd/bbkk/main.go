package main

import (
	"sync"

	"github.com/busyLambda/bbkk/internal/server"
)

func main() {
  var wg sync.WaitGroup
  wg.Add(1)

	mcServer := server.NewMcServer("server", "paper.jar", "")

	mcServer.Start(&wg)
  
  wg.Wait()
}
