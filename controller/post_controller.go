package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/keiya01/eserver/http/request"
	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service"
)

type PostController struct {
	*service.Service
}

func NewPostController() *PostController {
	return &PostController{}
}

func (p PostController) Index(w http.ResponseWriter, r *http.Request) {
	posts := []model.Post{}

	if err := p.FindAll(&posts, "created_at desc"); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(posts)
}

func (p PostController) Show(w http.ResponseWriter, r *http.Request) {
	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	post := model.Post{}
	if err := p.FindOne(&post, id); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(post)
}

func (p PostController) Create(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		panic(err)
	}

	if err := p.Service.Create(&post); err != nil {
		log.Printf("Create in PostController(Create()): %v", err)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (p PostController) Update(w http.ResponseWriter, r *http.Request) {
	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		panic(err)
	}
	params := map[string]interface{}{
		"name": post.Name,
		"body": post.Body,
		"url":  post.URL,
	}

	post.ID = id
	if err := p.Service.Update(&post, params); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(post)
}

func (p PostController) Delete(w http.ResponseWriter, r *http.Request) {
	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	var post model.Post
	if err := p.Service.Delete(&post, id); err != nil {
		panic(err)
	}

	message := struct {
		status  int
		message string
	}{
		status:  200,
		message: "削除に成功しました",
	}

	json.NewEncoder(w).Encode(message)
}
