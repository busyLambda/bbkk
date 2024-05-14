package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/busyLambda/bbkk/internal/models"
	"github.com/busyLambda/bbkk/internal/server"
	"github.com/busyLambda/bbkk/internal/util"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

func (a *App) startServer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	sid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := a.sm.GetServer(uint(sid))

	var wg sync.WaitGroup

	wg.Add(1)

	go s.Start(&wg)

	w.WriteHeader(http.StatusOK)
}

func (a *App) stopServer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	sid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := a.sm.GetServer(uint(sid))

	s.StopServer()
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
	q := chi.URLParam(r, "query")

	s, err := a.db.GetServerByName(q)
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

// func getServerCount(w http.ResponseWriter, r *http.Request) {
// 	q := chi.URLParam(r, "query")
// }

func (a *App) statusReport(w http.ResponseWriter, r *http.Request) {
	log.Println("WE GET HERE.")
	id := chi.URLParam(r, "id")
	sid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := a.sm.GetServer(uint(sid))

	// Use gorilla to get a websocket connection.
	c, err := a.up.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Success... so far anyway.")

	go statusReportStream(c, s)
}

func statusReportStream(c *websocket.Conn, s *server.McServer) {
	defer c.Close()

	// Used to decide if we need to report if the server is running.
	// report_r := true

	go func() {
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Client disconnected")
				} else {
					log.Println("Error reading message:", err)
				}
				return
			}
		}
	}()

	for {
		time.Sleep(time.Second)

		if !s.IsRunning() {
			c.WriteMessage(websocket.TextMessage, []byte(`{"not_running": true}`))
			continue
		}

		memuse, err := util.GetRssByPid(s.Cmd.Process.Pid)
		if err != nil {
			log.Println(err)
			c.WriteMessage(websocket.TextMessage, []byte(`{"failure": true}`))
			continue
		}

		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"is_running": true, "mem_use": %d}`, memuse)))
	}
}

func (a *App) openConsole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := a.sm.GetServer(uint(sid))

	c, err := a.up.Upgrade(w, r, nil)
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
