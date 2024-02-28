package api

import (
	"github.com/busyLambda/bbkk/internal/db"
	"github.com/go-chi/chi"
)

type App struct {
	h  *chi.Mux
	db *db.DbManager
}

func NewApiMaster() App {
	return App{
		h: chi.NewRouter(),
	}
}

func (a *App) AttachRoutes() {

}
