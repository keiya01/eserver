package controller

import (
	"encoding/json"
	"net/http"

	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service"
	"github.com/keiya01/eserver/service/database"
)

type PostController struct {
	*service.Service
}

func NewPostController() *PostController {
	return &PostController{}
}

func (p PostController) Index(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	service := service.NewService(db)
	posts := []model.Post{}
	service.FindAll(&posts, "created_at desc")

	json.NewEncoder(w).Encode(posts)
}
