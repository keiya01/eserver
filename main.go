package main

import (
	"github.com/keiya01/eserver/http"
	"github.com/keiya01/eserver/service/database"
	"github.com/keiya01/eserver/service/migrate"
)

func main() {
	DBHandler := database.NewHandler()
	migrate.Set(DBHandler.DB)
	s := http.NewServer()
	s.Router()
	s.Start(":8686")
}
