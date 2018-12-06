package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/keiya01/eserver/auth"
	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service"
	"github.com/keiya01/eserver/service/database"
	"github.com/keiya01/eserver/service/migrate"
	"github.com/stretchr/testify/assert"
)

var loc, _ = time.LoadLocation("Asia/Tokyo")

func TestMain(m *testing.M) {
	// サーバーをテスト用に変更
	os.Setenv("ENV", "TEST")

	test := m.Run()

	os.Exit(test)
}

func newMockModel() []interface{} {
	pass1, _ := auth.EncryptPassword("password")
	pass2, _ := auth.EncryptPassword("pass")
	pass3, _ := auth.EncryptPassword("word")
	return []interface{}{
		&model.User{
			Email:    "mail@mail.com",
			Password: pass1,
			Model: model.Model{
				CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				UpdatedAt: time.Time{},
			},
		},
		&model.User{
			Email:    "gmail@mail.com",
			Password: pass2,
			Model: model.Model{
				CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				UpdatedAt: time.Time{},
			},
		},
		&model.User{
			Email:    "mymail@mail.com",
			Password: pass3,
			Model: model.Model{
				CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				UpdatedAt: time.Time{},
			},
		},
		&model.Post{
			Name:   "Google",
			Body:   "aaaa",
			URL:    "https://www.google.com",
			UserID: 1,
			Model: model.Model{
				CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				UpdatedAt: time.Time{},
			},
		},
		&model.Post{
			Name:   "Yahoo",
			Body:   "aaaaaa",
			URL:    "https://www.yahoo.ne.jp",
			UserID: 2,
			Model: model.Model{
				CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				UpdatedAt: time.Time{},
			},
		},
		&model.Post{
			Name:   "Go",
			Body:   "aaaaaaa",
			URL:    "https://www.golang.org",
			UserID: 3,
			Model: model.Model{
				CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				UpdatedAt: time.Time{},
			},
		},
	}

}

func newMockServer() *service.Service {
	mock := newMockModel()
	db := database.NewHandler()
	service := service.NewService(db)
	defer service.Close()
	migrate.Set(db.DB)

	// モックデータを作成する
	for i := 0; i < len(mock); i++ {
		err := service.Create(mock[i])
		if err != nil {
			fmt.Printf("newMockModel: %v", err)
			return service
		}
	}

	return service
}

func newMockJWT(s *service.Service) string {
	// var user model.User
	token := auth.JWTToken{
		UserID:    1,
		UserEmail: "mail@mail.com",
	}
	jwt := token.GetJWTToken()

	return jwt
}

var client = new(http.Client)

func Test指定したパスにアクセスしたときにJSONを返すことを確認するテスト(t *testing.T) {
	type args struct {
		path    string
		method  string
		request io.Reader
		status  int
	}
	tests := []struct {
		name    string
		args    args
		want    model.Response
		wantErr bool
		hasJWT  bool
	}{
		{
			name: "/api/posts/{id}にアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/posts/1",
				method:  "GET",
				request: nil,
				status:  200,
			},
			want: model.Response{
				Data: model.Post{
					Name: "Google",
					Body: "aaaa",
					URL:  "https://www.google.com",
					Model: model.Model{
						ID:        0,
						CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						UpdatedAt: time.Time{},
					},
				},
			},
			hasJWT: true,
		},
		{
			name: "/api/posts/{id}にアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する",
			args: args{
				path:    "/api/posts/1",
				method:  "GET",
				request: nil,
				status:  403,
			},
			want: model.Response{
				Error: model.Error{
					IsErr:   true,
					Message: "認証エラーが発生しました",
				},
			},
		},
		{
			name: "/api/posts/createにアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/posts/create",
				method:  "POST",
				request: strings.NewReader(`{"name":"Hello","body":"bbbbb","url":"https://www.cash.com","created_at":"2014-12-31T08:04:18+09:00","updated_at":"2014-12-31T08:04:18+09:00"}`),
				status:  200,
			},
			want: model.Response{
				Message: "データを保存しました",
				Data: map[string]interface{}{
					"name":       "Hello",
					"body":       "bbbbb",
					"url":        "https://www.cash.com",
					"created_at": time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				},
			},
			hasJWT: true,
		},
		{
			name: "/api/posts/createにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する",
			args: args{
				path:    "/api/posts/create",
				method:  "POST",
				request: strings.NewReader(`{"name":"Hello","body":"bbbbb","url":"https://www.cash.com","created_at":"2014-12-31T08:04:18+09:00","updated_at":"2014-12-31T08:04:18+09:00"}`),
				status:  403,
			},
			want: model.Response{
				Error: model.Error{
					IsErr:   true,
					Message: "認証エラーが発生しました",
				},
			},
		},
		{
			name: "/api/posts/{id}/updateにアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/posts/1/update",
				method:  "PUT",
				request: strings.NewReader(`{"name":"Hello","body":"bbbbb","url":"https://www.cash.com","created_at":"2014-12-31T08:04:18+09:00","updated_at":"2014-12-31T08:04:18+09:00"}`),
				status:  200,
			},
			want: model.Response{
				Message: "データを更新しました",
				Data: map[string]interface{}{
					"name":       "Hello",
					"body":       "bbbbb",
					"url":        "https://www.cash.com",
					"created_at": time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				},
			},
			hasJWT: true,
		},
		{
			name: "/api/posts/{id}/updateにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する",
			args: args{
				path:    "/api/posts/1/update",
				method:  "PUT",
				request: strings.NewReader(`{"name":"Hello","body":"bbbbb","url":"https://www.cash.com","created_at":"2014-12-31T08:04:18+09:00","updated_at":"2014-12-31T08:04:18+09:00"}`),
				status:  403,
			},
			want: model.Response{
				Error: model.Error{
					IsErr:   true,
					Message: "認証エラーが発生しました",
				},
			},
		},
		{
			name: "/api/posts/{id}/deleteにアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/posts/1/delete",
				method:  "DELETE",
				request: nil,
				status:  200,
			},
			want: model.Response{
				Message: "データを削除しました",
			},
			hasJWT: true,
		},
		{
			name: "/api/posts/{id}/deleteにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する",
			args: args{
				path:    "/api/posts/1/delete",
				method:  "DELETE",
				request: nil,
				status:  403,
			},
			want: model.Response{
				Error: model.Error{
					IsErr:   true,
					Message: "認証エラーが発生しました",
				},
			},
		},
		{
			name: "/api/users/{id}にアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/users/1",
				method:  "GET",
				request: nil,
				status:  200,
			},
			want: model.Response{
				Message: "データを取得しました",
				Data: []model.Post{
					model.Post{
						Name: "Google",
						Body: "aaaa",
						URL:  "https://www.google.com",
						Model: model.Model{
							ID:        0,
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Time{},
						},
					},
				},
			},
			hasJWT: true,
		},
		{
			name: "/api/users/{id}にアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する",
			args: args{
				path:    "/api/users/1",
				method:  "GET",
				request: nil,
				status:  403,
			},
			want: model.Response{
				Error: model.Error{
					IsErr:   true,
					Message: "認証エラーが発生しました",
				},
			},
		},
		{
			name: "/api/users/{id}/updateにアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/users/1/update",
				method:  "PUT",
				request: strings.NewReader(`{"email":"testMail@test.mail.com","password":"testPassword","created_at":"2014-12-31T08:04:18+09:00","updated_at":"2014-12-31T08:04:18+09:00"}`),
				status:  200,
			},
			want: model.Response{
				Message: "データを更新しました",
				Data: map[string]interface{}{
					"email": "testMail@test.mail.com",
				},
			},
			hasJWT: true,
		},
		{
			name: "/api/users/{id}/updateにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する",
			args: args{
				path:    "/api/users/1/update",
				method:  "PUT",
				request: strings.NewReader(`{"email":"testMail@test.mail.com","password":"testPassword","created_at":"2014-12-31T08:04:18+09:00","updated_at":"2014-12-31T08:04:18+09:00"}`),
				status:  403,
			},
			want: model.Response{
				Error: model.Error{
					IsErr:   true,
					Message: "認証エラーが発生しました",
				},
			},
		},
		{
			name: "/api/users/{id}/deleteにアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/users/1/delete",
				method:  "DELETE",
				request: nil,
				status:  200,
			},
			want: model.Response{
				Message: "データを削除しました",
			},
			hasJWT: true,
		},
		{
			name: "/api/users/{id}/deleteにアクセスしたときにヘッダーのAuthorizationにJWTを持っていなければエラーを返すことを確認する",
			args: args{
				path:    "/api/users/1/delete",
				method:  "DELETE",
				request: nil,
				status:  403,
			},
			want: model.Response{
				Error: model.Error{
					IsErr:   true,
					Message: "認証エラーが発生しました",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// DBのデータモック作成用関数
			service := newMockServer()
			defer os.Remove("test.sqlite3")
			token := newMockJWT(service)

			server := NewServer()
			testServer := httptest.NewServer(server)
			server.Router()
			defer testServer.Close()

			// リクエスト先のURLを作成し、リクエストを送る
			path := testServer.URL + tt.args.path
			req, _ := http.NewRequest(tt.args.method, path, tt.args.request)
			req.Header.Set("Content-Type", "application/json")

			if tt.hasJWT {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("client.Do(): %v", err)
			}
			defer resp.Body.Close()

			// レスポンスで受け取ったJSONを文字列に変換する
			respBody, _ := ioutil.ReadAll(resp.Body)
			respJSON := string(respBody)

			// 期待する値を文字列のJSONへ変換する
			wantJSON, _ := json.Marshal(tt.want)
			want := string(wantJSON)

			assert.Equal(t, tt.args.status, resp.StatusCode)

			assert.JSONEq(t, want, respJSON)
		})
	}
}

func TestUserControllerのcreateまたはloginにアクセスしたときに正しいJSONを返すことを確認する(t *testing.T) {
	type args struct {
		path    string
		request io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "/api/users/createにアクセスしたときにJSONを返すことを確認すること",
			args: args{
				path:    "/api/users/create",
				request: strings.NewReader(`{"email":"mailmail@mailmail.co.jp","password":"passtestpass","created_at":"2014-12-31T08:04:18+09:00","updated_at":"2014-12-31T08:04:18+09:00"}`),
			},
			want: "データを保存しました",
		},
		{
			name: "/api/users/loginにアクセスしたときにJSONを返すことを確認すること",
			args: args{
				path:    "/api/users/login",
				request: strings.NewReader(`{"email":"mail@mail.com","password":"password"}`),
			},
			want: "ログインに成功しました",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// DBのデータモック作成用関数
			_ = newMockServer()
			defer os.Remove("test.sqlite3")

			server := NewServer()
			testServer := httptest.NewServer(server)
			server.Router()
			defer testServer.Close()

			// リクエスト先のURLを作成し、リクエストを送る
			path := testServer.URL + tt.args.path
			req, _ := http.NewRequest("POST", path, tt.args.request)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("client.Do(): %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)

		})
	}
}
