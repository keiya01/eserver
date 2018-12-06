package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/keiya01/eserver/auth"
	"github.com/keiya01/eserver/http/request"
	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service"
	"github.com/keiya01/eserver/service/database"
)

type UserController struct{}

func (u UserController) Login(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	var params model.User
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		panic(err)
	}

	var user model.User
	var response model.Response
	if err := s.FindOne(&user, "email = ?", params.Email); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		response.Error = model.NewError("メールアドレスが見つかりませんでした")

		json.NewEncoder(w).Encode(response)

		return
	}

	if auth.ComparePassword(user.Password, params.Password) {
		//Create JWT token
		tk := &auth.JWTToken{UserID: user.ID, UserEmail: user.Email}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
		response = model.Response{
			Token:   tokenString, //Store the token in the response
			Message: "ログインに成功しました",
			Data:    user,
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		response.Error = model.NewError("パスワードが一致しませんでした")
	}

	json.NewEncoder(w).Encode(response)
}
func (u UserController) Create(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	log.Printf("create: %v", r.Body)
	var params model.User
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		panic(err)
	}

	encryptedPassword, err := auth.EncryptPassword(params.Password)
	if err != nil {
		panic(err)
	}

	user := model.User{
		Email:    params.Email,
		Password: encryptedPassword,
	}

	var resp model.Response
	if err := s.Create(&user); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを保存できませんでした")

		json.NewEncoder(w).Encode(resp)
		return
	}

	token := auth.JWTToken{
		UserID:    user.ID,
		UserEmail: user.Email,
	}
	jwtToken := token.GetJWTToken()

	resp.Token = jwtToken
	resp.Message = "データを保存しました"

	json.NewEncoder(w).Encode(resp)

}

func (u UserController) Show(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	// ページネーションを行うためのデータを取得
	queryVal := r.URL.Query()
	page := queryVal.Get("page")
	var pageNum int
	if page != "" {
		var err error
		pageNum, err = strconv.Atoi(page)
		if err != nil {
			panic(err)
		}
	}

	// 初期ページのデータを下にoffsetの数値を取得
	initPage := 20
	nextPage := initPage * pageNum

	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	var resp model.Response

	var posts []model.Post
	if err := s.Select("name, body, url, created_at").Pagination(initPage, nextPage).FindAll(&posts, "created_at desc", "user_id = ?", id); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを見つけることが出来ませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	resp = model.Response{
		Message: "データを取得しました",
		Data:    posts,
	}

	json.NewEncoder(w).Encode(resp)
}

func (u UserController) Update(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Printf("getting json data: %v", r.Body)
		panic(err)
	}

	encryptedPassword, err := auth.EncryptPassword(user.Password)
	if err != nil {
		panic(err)
	}
	params := map[string]interface{}{
		"email":    user.Email,
		"password": encryptedPassword,
	}

	var resp model.Response

	user.ID = id
	if err := s.Update(&user, params); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを更新できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	resp = model.Response{
		Message: "データを更新しました",
		Data: map[string]interface{}{
			"email": user.Email,
		},
	}

	json.NewEncoder(w).Encode(resp)

}

func (u UserController) Delete(w http.ResponseWriter, r *http.Request) {
	db := database.NewHandler()
	s := service.NewService(db)
	defer s.Close()

	var resp model.Response

	paramsID := request.GetParam(r, "id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		panic(err)
	}

	var user model.User
	if err := s.Delete(&user, id); err != nil {
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
