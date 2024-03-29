package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/busyLambda/bbkk/internal/db"
	"github.com/busyLambda/bbkk/internal/server"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type App struct {
	r  *chi.Mux
	db *db.DbManager
	sm *server.ServerManager
}

func NewApiMaster() App {
	db := db.NewDbManager("localhost", "unkindled", "CoCk1234", "bbkk_dev", 5432, "Europe/Budapest")

	log.Println("Connected to database, getting servers.")

	sm := server.NewServerManager()

	servers, err := db.GetAllServers()
	if err != nil {
		// TODO: This seems needless.
		if err == gorm.ErrRecordNotFound {
			log.Println("No servers found.")
		} else {
			log.Fatalf("Error getting servers: %s", err)
		}
	} else {
		log.Println("Servers found.")
		for _, s := range servers {
			sm.AddServer(s.ID, server.NewMcServer(s.Name, "", ""))
		}
	}

	return App{
		r:  chi.NewRouter(),
		db: db,
		sm: sm,
	}
}

func (a *App) AttachRoutes() {
	a.r.Post("/register", a.createUser)
}

func (a *App) Run(port uint) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), a.r)
}
