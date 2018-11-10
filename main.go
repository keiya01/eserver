package main

import (
	"github.com/keiya01/eserver/database"
	"github.com/keiya01/eserver/database/migrate"
	"github.com/keiya01/eserver/http"
)

func main() {
	DBHandler := database.NewHandler()
	migrate.Set(DBHandler.DB)
	s := http.NewServer()
	s.Router()
	s.Start("8080")
}
