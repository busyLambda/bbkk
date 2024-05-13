package api

import (
	"context"
	"encoding/json"
	"fmt"
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

// TODO: Use a token with a key that we can use to validate the session and also still check the user in the DB.
// TODO: return user.ID and user.Username as JSON.
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

	cookie := http.Cookie{
		Domain:   "localhost",
		Name:     "session",
		Value:    s.ID,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

  resp := map[string]string{
    "id": fmt.Sprint(u.ID),
    "username": u.Username,
  }

	w.Header().Set("Content-Type", "application/json")
  err = json.NewEncoder(w).Encode(resp)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

// TODO: Don't respond with the whole user.
func (a *App) validateSession(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(util.UserKey{})

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
