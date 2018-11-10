package service

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service/database"
	"github.com/keiya01/eserver/service/migrate"
)

func TestMain(m *testing.M) {
	handler := database.NewTestHandler()
	defer handler.DB.Close()
	migrate.Set(handler.DB)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	post := model.Post{
		Title: "Hello",
		Body:  "World",
		Model: model.Model{
			CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
			UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
		},
	}
	handler.DB.Create(&post)

	code := m.Run()
	defer os.Exit(code)

}

func Testデータベースに保存されている投稿データを全て取得することを確認する(t *testing.T) {
	time := `"2014-12-31T08:04:18+09:00"`
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "3つの投稿データを取得できていることを確認する",
			want: `{"id":1,"created_at":` + time + `,"updated_at":` + time + `,"title":"Hello","body":"World"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := database.NewTestHandler()
			p := PostService{
				handler: handler,
			}
			posts, _ := p.FindAll()
			json, _ := json.Marshal(posts)
			get := string(json)
			if get != tt.want {
				t.Errorf("FindAll() get = %v want = %v", get, tt.want)
			}
		})
	}
}
