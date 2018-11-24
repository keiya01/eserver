package controller

import (
	"encoding/json"
	"fmt"
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

func (p PostController) Show(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	post := model.Post{}
	if err := s.Select("name, body, url, created_at").FindOne(&post, id); err != nil {
		panic(err)
	}

	resp := model.Response{
		Status: 200,
		Data:   post,
	}

	json.NewEncoder(w).Encode(resp)
}

func (p PostController) Create(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		panic(err)
	}

	if err := s.Create(&post); err != nil {
		log.Printf("Create in PostController(Create()): %v", err)
		return
	}

	resp := model.Response{
		Status:  200,
		Message: "データを保存しました",
	}

	json.NewEncoder(w).Encode(resp)
}

func (p PostController) Update(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		fmt.Printf("getting json data: %v", r.Body)
		panic(err)
	}
	params := map[string]interface{}{
		"name": post.Name,
		"body": post.Body,
		"url":  post.URL,
	}

	post.ID = id
	if err := s.Update(&post, params); err != nil {
		panic(err)
	}

	resp := model.Response{
		Status:  200,
		Message: "データを更新しました",
	}

	json.NewEncoder(w).Encode(resp)
}

func (p PostController) Delete(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	var post model.Post
	if err := s.Delete(&post, id); err != nil {
		panic(err)
	}

	resp := model.Response{
		Status:  200,
		Message: "データを削除しました",
	}

	json.NewEncoder(w).Encode(resp)
}
