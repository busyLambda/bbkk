package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/busyLambda/bbkk/domain/user"
	"github.com/busyLambda/bbkk/internal/db"
	"github.com/busyLambda/bbkk/internal/models"
	"github.com/busyLambda/bbkk/internal/server"
	"github.com/busyLambda/bbkk/internal/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	r  *chi.Mux
	db *db.DbManager
	sm *server.ServerManager
	up *websocket.Upgrader
}

func NewApiMaster() App {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	host := os.Getenv("BBKK_DB_HOST")
	port_var := os.Getenv("BBKK_DB_PORT")

	port, err := strconv.ParseUint(port_var, 10, 32)
	if err != nil {
		log.Fatalf("Error parsing port: %s from env var %s", err, port_var)
	}

	username := os.Getenv("BBKK_DB_USER")
	dbname := os.Getenv("BBKK_DB_NAME")
	pass := os.Getenv("BBKK_DB_PASS")
	locale := os.Getenv("BBKK_DB_LOCALE")

	db := db.NewDbManager(host, username, pass, dbname, uint(port), locale)

	log.Println("Connected to database, getting servers.")

	var superAdmin models.User
	err = db.Conn.Where("role = ?", user.SUPERADMIN).First(&superAdmin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("No super admin found, creating one.\n\n")

			var username string
			for {
				fmt.Printf("Enter a username: ")
				fmt.Scanln(&username)
				if len(username) < 1 {
					fmt.Println("Username must be more than 1 character.")
				} else if len(username) > 48 {
					fmt.Println("Username must be less than 48 characters.")
				} else {
					break
				}
			}

			var password string
			for {
				fmt.Printf("Enter a password (must be more than 12 characters): ")
				fmt.Scanln(&password)
				if len(password) < 12 {
					fmt.Println("Password must be more than 12 characters.")
				} else {
					break
				}
			}

			superAdmin, err = models.NewUser(util.RegistrationForm{Username: username, Password: password}, user.SUPERADMIN)
			if err != nil {
				log.Fatalf("Error creating super admin: %s", err)
			}

			err = db.InsertUser(&superAdmin)
			if err != nil {
				log.Fatalf("Error inserting super admin: %s", err)
			}

			fmt.Println("Super admin created.")
		} else {
			log.Fatalf("Error getting super admin: %s", err)
		}
	} else {
		log.Printf("Super admin: %s\n", superAdmin.Username)
	}

	sm := server.NewServerManager()

	servers, err := db.GetAllServers()
	if err != nil {
		// TODO: This seems needless.
		if err == gorm.ErrRecordNotFound {
			log.Println("No servers found.")
		} else {
			log.Fatalf("Error getting servers: %s", err)
		}
	} else {
		log.Println("Finished getting servers.")
		i := 1
		for _, s := range servers {
			log.Printf("[%d]: %s\n", i, s.Name)
			sm.AddServer(s.ID, server.NewMcServer(util.ServerDirName(s.Name, s.ID), "server.jar", ""))
			i++
		}
	}
	upgrader := websocket.Upgrader{
		ReadBufferSize:  256,
		WriteBufferSize: 256,
		WriteBufferPool: &sync.Pool{},
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return App{
		r:  chi.NewRouter(),
		db: db,
		sm: sm,
		up: &upgrader,
	}
}

func (a *App) AttachRoutes() {
	a.r.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"https://localhost:5173", "http://localhost:5173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		},
	))

	// TODO: Make it one middleware `Use()`
	serverApi := chi.NewRouter()
	serverApi.Use(a.authMiddleware)
	userApi := chi.NewRouter()
	userApi.Use(a.authMiddleware)
	accessApi := chi.NewRouter()
	accessApi.Use(a.authMiddleware)

	a.r.Post("/login", a.login)                           // POST /login
	serverApi.Post("/console/write/{id}", a.writeConsole) // POST /server/console/write/{id}
	serverApi.Get("/console/{id}", a.openConsole)         // GET  /server/console/{id}
	serverApi.Post("/create", a.createServer)             // POST /server/create
	serverApi.Get("/name/{query}", a.getServerByName)     // GET  /server/name/{query}
	serverApi.Get("/all", a.getAllServers)                // GET  /server/all
	serverApi.Get("/start/{id}", a.startServer)           // GET  /server/start/{id}
	serverApi.Get("/statusreport/{id}", a.statusReport)   // GET  /server/statusreport/{id}
	serverApi.Get("/stop/{id}", a.stopServer)             // GET  /server/stop/{id}

	accessApi.Get("/validate", a.validateSession) // GET  /validate

	userApi.Get("/create", a.createUser) // POST /user/create

	a.r.Mount("/", accessApi)       // Protected
	a.r.Mount("/user", userApi)     // Protected
	a.r.Mount("/server", serverApi) // Protected
}

func (a *App) Run(port uint) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), a.r)
}
