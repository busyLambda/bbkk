package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/busyLambda/bbkk/internal/db"
	"github.com/busyLambda/bbkk/internal/server"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	r  *chi.Mux
	db *db.DbManager
	sm *server.ServerManager
}

func NewApiMaster() App {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	host := os.Getenv("BBKK_DB_HOST")
	port_var := os.Getenv("BBKK_DB_PORT")

	port, err := strconv.ParseUint(port_var, 10, 32)
	if err != nil {
		log.Fatalf("Error parsing port: %s from env var %s", err, port_var)
	}

	username := os.Getenv("BBKK_DB_USER")
	dbname := os.Getenv("BBKK_DB_NAME")
	pass := os.Getenv("BBKK_DB_PASS")
	locale := os.Getenv("BBKK_DB_LOCALE")

	db := db.NewDbManager(host, username, pass, dbname, uint(port), locale)

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
