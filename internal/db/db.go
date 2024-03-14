package db

import (
	"fmt"

	"github.com/busyLambda/bbkk/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbManager struct {
	Conn *gorm.DB
}

func NewDbManager(host string, user string, password string, dbName string, port uint, timeZone string) *DbManager {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", host, user, password, dbName, port, timeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Cannot open DB!")
	}

	err = db.AutoMigrate(models.User{})
	if err != nil {
		panic("Cannot migrate `models.User`")
	}

	return &DbManager{Conn: db}
}
