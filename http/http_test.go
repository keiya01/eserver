package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service"
	"github.com/keiya01/eserver/service/database"
	"github.com/keiya01/eserver/service/migrate"
)

func newMockModel() {
	mock := []interface{}{
		&model.Post{
			Name: "Google",
			URL:  "https://www.google.com",
		},
		&model.Post{
			Name: "Yahoo",
			URL:  "https://www.yahoo.ne.jp",
		},
		&model.Post{
			Name: "Go",
			URL:  "https://www.golang.org",
		},
	}

	db := database.NewHandler()
	service := service.NewService(db)
	defer service.Close()
	migrate.Set(db.DB)

	for i := 0; i < len(mock); i++ {
		err := service.Create(mock[i])
		if err != nil {
			fmt.Printf("newMockModel: %v", err)
			return
		}
	}
}

var client = new(http.Client)

func Test指定したパスにアクセスしたときにHTTPstatus200を返すことを確認するテスト(t *testing.T) {
	bodyReader := strings.NewReader(`{"name": "Hello", "url": "https://www.cash.com"}`)
	type args struct {
		path    string
		method  string
		request io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "/api/postsにアクセスしたときにHTTPstatusの200を返すことを確認する",
			args: args{
				path:    "/api/posts",
				method:  "GET",
				request: nil,
			},
			want: 200,
		},
		{
			name: "/api/posts/{id}にアクセスしたときにHTTPstatusの200を返すことを確認する",
			args: args{
				path:    "/api/posts/1",
				method:  "GET",
				request: nil,
			},
			want: 200,
		},
		{
			name: "/api/posts/createにアクセスしたときにHTTPstatusの200を返すことを確認する",
			args: args{
				path:    "/api/posts/create",
				method:  "POST",
				request: bodyReader,
			},
			want: 200,
		},
		{
			name: "/api/posts/{id}/updateにアクセスしたときにHTTPstatusの200を返すことを確認する",
			args: args{
				path:    "/api/posts/1/update",
				method:  "PUT",
				request: bodyReader,
			},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// DBのデータモック作成用関数
			newMockModel()
			defer os.Remove("eserver.sqlite3")

			server := NewServer()
			testServer := httptest.NewServer(server)

			server.Router()
			defer testServer.Close()

			path := testServer.URL + tt.args.path
			req, _ := http.NewRequest(tt.args.method, path, tt.args.request)
			resp, _ := client.Do(req)

			if resp.StatusCode != tt.want {
				t.Errorf("HTTP Test : get = %d", resp.StatusCode)
			}
		})
	}
}
