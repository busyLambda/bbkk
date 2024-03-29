package main

import (
	"log"

	"github.com/busyLambda/bbkk/internal/api"
)

func main() {
	app := api.NewApiMaster()

	app.AttachRoutes()

	port := 3000

	log.Printf("Server is running on port %d", port)
	app.Run(3000)
}
