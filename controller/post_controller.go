package controller

import (
	"encoding/json"
	"net/http"

	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service"
)

type PostController struct {
	PostService model.PostService
}

func NewPostController() *PostController {
	return &PostController{}
}

func (p PostController) Index(w http.ResponseWriter, r *http.Request) {
	postService := service.NewPostService()
	p.PostService = postService
	posts, err := p.PostService.FindAll()
	if err != nil {
		err := model.NewError(err)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(posts)
	}
}
