package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbManager struct {
	Conn *gorm.DB
}

func NewDbManager(host string, user string, password string, dbName string, port uint, timeZone string) DbManager {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", host, user, password, dbName, password, timeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Cannot open DB!")
	}

	return DbManager{Conn: db}
}
