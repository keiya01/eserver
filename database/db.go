package database

import (
	"github.com/jinzhu/gorm"
)

type DBHandler struct {
	DB *gorm.DB
}

func NewDBHandler() *DBHandler {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic(err)
	}

	handler := DBHandler{
		DB: db,
	}

	return &handler
}
