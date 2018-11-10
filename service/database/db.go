package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler() *Handler {
	db, err := gorm.Open("sqlite3", "blog.sqlite3")
	if err != nil {
		panic(err)
	}

	handler := Handler{
		DB: db,
	}

	return &handler
}

func NewTestHandler() *Handler {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic(err)
	}

	handler := Handler{
		DB: db,
	}

	return &handler
}
