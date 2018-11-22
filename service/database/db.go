package database

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Handler struct {
	*gorm.DB
}

func NewHandler() *Handler {
	isTest := os.Getenv("ENV")
	dbName := "eserver.sqlite3"
	if isTest == "TEST" {
		dbName = "test.sqlite3"
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	handler := Handler{db}

	return &handler
}
