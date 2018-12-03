package controller

import (
	"encoding/json"
	"fmt"
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

	var resp model.Response

	post := model.Post{}
	err = s.Select("name, body, url, created_at").FindOne(&post, id)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを取得できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	resp.Data = post

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

	var resp model.Response
	ctx := r.Context().Value("userID")
	userID, ok := ctx.(float64)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを取得できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	post.UserID = int(userID)

	if err := s.Create(&post); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを保存できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	resp = model.Response{
		Message: "データを保存しました",
		Data: map[string]interface{}{
			"name":       post.Name,
			"body":       post.Body,
			"url":        post.URL,
			"created_at": post.CreatedAt,
		},
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

	var resp model.Response

	post.ID = id
	if err := s.Update(&post, params); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを更新できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	resp = model.Response{
		Message: "データを更新しました",
		Data: map[string]interface{}{
			"name":       post.Name,
			"body":       post.Body,
			"url":        post.URL,
			"created_at": post.CreatedAt,
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func (p PostController) Delete(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	var resp model.Response

	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	var post model.Post
	if err := s.Delete(&post, id); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを削除できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	resp = model.Response{
		Message: "データを削除しました",
	}

	json.NewEncoder(w).Encode(resp)
}
