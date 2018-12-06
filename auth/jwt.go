package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/keiya01/eserver/model"
)

type JWTToken struct {
	UserID    int    `json:"userID"`
	UserEmail string `json:"userEmail"`
	jwt.StandardClaims
}

// GetTokenHandler get token
func (j *JWTToken) GetJWTToken() string {

	// headerのセット
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": true,
		"sub":   "http://localhost:8686/",
		"aud": []string{
			"http://localhost:8686/",
			"http://localhost:3000",
		},
		"iat":       time.Now(),
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
		"userEmail": j.UserEmail,
		"userID":    j.UserID,
	})

	// 電子署名
	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

	return tokenString
}

//ここでログインしているかどうかの検証をしている
func JWTAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/users/create", "/api/users/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                    //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		var resp model.Response
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
		if tokenHeader == "" {                       //Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			resp.Error = model.NewError("認証エラーが発生しました")

			json.NewEncoder(w).Encode(resp)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			resp.Error = model.NewError("ヘッダーの形式が正しくありません")

			json.NewEncoder(w).Encode(resp)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in

		token, err := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SIGNINGKEY")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			log.Printf("jwt.Parse: %v", err)
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			resp.Error = model.NewError("トークンを認証できませんでした")

			json.NewEncoder(w).Encode(resp)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			resp.Error = model.NewError("トークンが認証されませんでした")

			json.NewEncoder(w).Encode(resp)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			resp.Error = model.NewError("トークンが認証されませんでした")

			json.NewEncoder(w).Encode(resp)
			return
		}

		log.Printf("JWT: %v", claims["userID"])

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), "userID", claims["userID"]) // "userID"にclaims[UserId]を入れる
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
