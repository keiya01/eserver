package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/keiya01/eserver/http/request"
	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service"
	"github.com/keiya01/eserver/service/database"
)

type PostController struct {
}

func NewPostController() *PostController {
	return &PostController{}
}

func (p PostController) Index(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	posts := []model.Post{}

	if err := s.FindAll(&posts, "created_at desc"); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(posts)
}

func (p PostController) Show(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)

	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	post := model.Post{}
	if err := s.FindOne(&post, id); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(post)
}

func (p PostController) Create(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)

	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		panic(err)
	}

	if err := s.Create(&post); err != nil {
		log.Printf("Create in PostController(Create()): %v", err)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (p PostController) Update(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)

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
		"url":  post.URL,
	}

	post.ID = id
	if err := s.Update(&post, params); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(post)
}
