package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/busyLambda/bbkk/internal/models"
	"github.com/busyLambda/bbkk/internal/util"
)

func (a *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Unauthorized: no session found.", http.StatusUnauthorized)
			return
		}

		s, err := a.db.GetSessionById(c.Value)
		if err != nil {
			http.Error(w, "Unauthorized: session not found.", http.StatusUnauthorized)
			return
		}

		u, err := a.db.GetUserByID(int(s.UserID))
		if err != nil {
			http.Error(w, "Unauthorized: no such user found.", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), util.UserKey{}, u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var login util.LoginForm

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := a.db.GetUserByUsername(login.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !util.CheckPasswordHash(login.Password, u.Password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	s := models.NewSession(u.ID, r.UserAgent())

	err = a.db.InsertSession(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(s.ID))
}
