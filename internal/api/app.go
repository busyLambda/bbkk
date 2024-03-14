package api

import (
	"net/http"

	"github.com/busyLambda/bbkk/internal/db"
	"github.com/go-chi/chi"
)

type App struct {
	r  *chi.Mux
	db *db.DbManager
}

func NewApiMaster() App {
	db := db.NewDbManager("localhost", "unkindled", "CoCk1234", "bbkk_dev", 5432, "Europe/Budapest")

	return App{
		r:  chi.NewRouter(),
		db: db,
	}
}

func (a *App) AttachRoutes() {
	a.r.Post("/register", a.createUser)
}

func (a *App) Run() {
	http.ListenAndServe(":3000", a.r)
}
