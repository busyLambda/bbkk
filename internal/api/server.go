package api

import (
	"encoding/json"
	"net/http"

	"github.com/busyLambda/bbkk/internal/models"
	"github.com/go-chi/chi"
)

func (a *App) createServer(w http.ResponseWriter, r *http.Request) {
	var s models.Server

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.db.InsertServer(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getServerByName(w http.ResponseWriter, r *http.Request) {
	n := chi.URLParam(r, "name")

	s, err := a.db.GetServerByName(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}
