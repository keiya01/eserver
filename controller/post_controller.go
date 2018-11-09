package controller

import (
	"encoding/json"
	"net/http"

	"github.com/keiya01/eserver/model"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (p PostController) Index(w http.ResponseWriter, r *http.Request) {
	posts := model.NewPost("Test desu", "Hello World")

	json.NewEncoder(w).Encode(posts)
}
