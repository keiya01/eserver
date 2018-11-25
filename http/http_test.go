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
	return []interface{}{
		&model.Post{
			Name: "Google",
			Body: "aaaa",
			URL:  "https://www.google.com",
			Model: model.Model{
				CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				UpdatedAt: time.Time{},
			},
		},
		&model.Post{
			Name: "Yahoo",
			Body: "aaaaaa",
			URL:  "https://www.yahoo.ne.jp",
			Model: model.Model{
				CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				UpdatedAt: time.Time{},
			},
		},
		&model.Post{
			Name: "Go",
			Body: "aaaaaaa",
			URL:  "https://www.golang.org",
			Model: model.Model{
				CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
				UpdatedAt: time.Time{},
			},
		},
	}

}

func newMockServer() {
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
			return
		}
	}
}

var client = new(http.Client)

func Test指定したパスにアクセスしたときにJSONを返すことを確認するテスト(t *testing.T) {
	type args struct {
		path    string
		method  string
		request io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    model.Response
		wantErr bool
	}{
		{
			name: "/api/posts/{id}にアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/posts/1",
				method:  "GET",
				request: nil,
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
		},
		{
			name: "/api/posts/createにアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/posts/create",
				method:  "POST",
				request: strings.NewReader(`{"name":"Hello","body":"bbbbb","url":"https://www.cash.com","created_at":"2014-12-31T08:04:18+09:00","updated_at":"2014-12-31T08:04:18+09:00"}`),
			},
			want: model.Response{
				Message: "データを保存しました",
			},
		},
		{
			name: "/api/posts/{id}/updateにアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/posts/1/update",
				method:  "PUT",
				request: strings.NewReader(`{"name":"Hello","body":"bbbbb","url":"https://www.cash.com","created_at":"2014-12-31T08:04:18+09:00","updated_at":"2014-12-31T08:04:18+09:00"}`),
			},
			want: model.Response{
				Message: "データを更新しました",
			},
		},
		{
			name: "/api/posts/{id}/deleteにアクセスしたときにJSONを返すことを確認する",
			args: args{
				path:    "/api/posts/1/delete",
				method:  "DELETE",
				request: nil,
			},
			want: model.Response{
				Message: "データを削除しました",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// DBのデータモック作成用関数
			newMockServer()
			defer os.Remove("test.sqlite3")

			server := NewServer()
			testServer := httptest.NewServer(server)
			server.Router()
			defer testServer.Close()

			// リクエスト先のURLを作成し、リクエストを送る
			path := testServer.URL + tt.args.path
			req, _ := http.NewRequest(tt.args.method, path, tt.args.request)
			req.Header.Set("Content-Type", "application/json")

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

			assert.Equal(t, 200, resp.StatusCode)

			assert.JSONEq(t, want, respJSON)
		})
	}
}
