package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/busyLambda/bbkk/internal/models"
	"github.com/busyLambda/bbkk/internal/server"
	"github.com/busyLambda/bbkk/internal/util"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  256,
	WriteBufferSize: 256,
	WriteBufferPool: &sync.Pool{},
}

func (a *App) createServer(w http.ResponseWriter, r *http.Request) {
	var sf util.ServerForm

	err := json.NewDecoder(r.Body).Decode(&sf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := models.NewServer(&sf)

	err = a.db.InsertServer(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.sm.AddServer(s.ID, server.NewMcServer(s.Name, "server.jar", ""))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (a *App) getServerByName(w http.ResponseWriter, r *http.Request) {
	n := chi.URLParam(r, "name")

	s, err := a.db.GetServerByName(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (a *App) getAllServers(w http.ResponseWriter, r *http.Request) {
	s, err := a.db.GetAllServers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (a *App) openConsole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := a.sm.GetServer(uint(sid))

	// Use gorilla to get a websocket connection.
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go console(c, s)
}

func console(c *websocket.Conn, s *server.McServer) {
	defer c.Close()

	out := make(chan string)

	if !s.IsStreaming() {
		s.SetStdout()
		s.SetStdin()
		s.ReadStdout(out)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		o := <-out
		c.WriteMessage(websocket.TextMessage, []byte(o))
	}
}
