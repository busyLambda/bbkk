package main

import (
	"github.com/busyLambda/bbkk/internal/server"
)

func main() {
	mcServer := server.NewMcServer("server", "paper.jar")

	mcServer.Start()
}
