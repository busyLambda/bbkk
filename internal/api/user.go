package api

import (
	"encoding/json"
	"net/http"

	"github.com/busyLambda/bbkk/domain/user"
	"github.com/busyLambda/bbkk/internal/models"
	"github.com/busyLambda/bbkk/internal/util"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var rf util.RegistrationForm

	err := json.NewDecoder(r.Body).Decode(&rf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := models.NewUser(rf, user.ADMIN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.db.InsertUser(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Success 200"))
}

func (a *App) getUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	u, err := a.db.GetUserByUsername(username)
	if err != nil {
		code := http.StatusInternalServerError

		if err == gorm.ErrRecordNotFound {
			code = http.StatusNotFound
		}

		http.Error(w, err.Error(), code)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := a.db.Conn.Delete(models.User{}, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
