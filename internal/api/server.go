package api

import (
	"encoding/json"
	"fmt"
	"io"
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

	// TODO: Handle this better if it fails.
	err = util.CreateServer(s.Name, s.ID, sf.Version, sf.Build)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.sm.AddServer(s.ID, server.NewMcServer(fmt.Sprintf("servers/%s-%d", s.Name, s.ID), "server.jar", ""))

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

func (a *App) statusReport(w http.ResponseWriter, r *http.Request) {
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

	go statusReportStream(c, s)
}

func statusReportStream(c *websocket.Conn, s *server.McServer) {
	defer c.Close()

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

	// TODO: Rename this, tv stands for: temporary variable, and tv_p is temoporary variable previous
	tv := true
	tv_p := false

	for {
		time.Sleep(time.Second)

		// Tbh I don't know the logic, was feeling fuzzy so I just kinda wrote something and it works so IDK.
		if tv != tv_p {
			if tv {
				tv = false
				if !s.IsRunning() {
					c.WriteMessage(websocket.TextMessage, []byte(`{"not_running": true}`))
					continue
				}

			}
		}

		if !s.IsRunning() {
			continue
		} else {
			if !tv {
				tv = true
			}
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
		if s.Cmd.ProcessState != nil {
			if s.Cmd.ProcessState.Exited() {
				c.Close()
				return
			}
		}

		for s.Stdout.Scan() {
			text := s.Stdout.Text()
			fmt.Println(text)
			c.WriteMessage(websocket.TextMessage, []byte(text))
			// time.Sleep(time.Millisecond * 50)
		}
	}
}

func (a *App) writeConsole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := a.sm.GetServer(uint(sid))

	if !s.IsRunning() {
		http.Error(w, "Server is not running.", http.StatusInternalServerError)
		return
	}

	// if !s.IsStreaming() {
	// 	http.Error(w, "Server is not streaming console.", http.StatusInternalServerError)
	// 	return
	// }

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body.", http.StatusInternalServerError)
		return
	}

	input := fmt.Sprintf("%s\n", string(bytes))

	s.WriteString(input)
}
