package api

import (
	"encoding/json"
	"net/http"

	"github.com/busyLambda/bbkk/internal/models"
	"github.com/busyLambda/bbkk/internal/util"
	"github.com/go-chi/chi"
)

func (a *App) AddUser(w http.ResponseWriter, r *http.Request) {
	var rf util.RegistrationForm

	err := json.NewDecoder(r.Body).Decode(&rf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := models.NewUser(rf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.db.InsertUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := a.db.Conn.Delete(models.User{}, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
