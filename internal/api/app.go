package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/busyLambda/bbkk/domain/user"
	"github.com/busyLambda/bbkk/internal/db"
	"github.com/busyLambda/bbkk/internal/models"
	"github.com/busyLambda/bbkk/internal/server"
	"github.com/busyLambda/bbkk/internal/util"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	r  *chi.Mux
	db *db.DbManager
	sm *server.ServerManager
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
			log.Printf("[%d.]: %s\n", i, s.Name)
			sm.AddServer(s.ID, server.NewMcServer(s.Name, "", ""))
			i++
		}
	}

	return App{
		r:  chi.NewRouter(),
		db: db,
		sm: sm,
	}
}

func (a *App) AttachRoutes() {
	a.r.Post("/register", a.createUser)
	a.r.Post("/login", a.login)

	serverApi := chi.NewRouter()

	serverApi.Use(a.authMiddleware)

	serverApi.Get("/console/{id}", a.openConsole)
	serverApi.Post("/create", a.createServer)

	a.r.Mount("/server", serverApi)
}

func (a *App) Run(port uint) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), a.r)
}
