package service

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service/database"
	"github.com/keiya01/eserver/service/migrate"
)

var loc, _ = time.LoadLocation("Asia/Tokyo")

func startMockDB(mockModel []interface{}, queryTest func(s *Service)) {
	handler := database.NewTestHandler()
	service := NewService(handler)
	defer os.Remove("test.sqlite3")
	defer handler.Close()
	migrate.Set(handler.DB)
	for i := 0; i < len(mockModel); i++ {
		if err := service.Create(mockModel[i]); err != nil {
			panic(err)
		}
	}

	queryTest(service)
}

func TestデータベースからIDに一致しているデータを取り出すことを確認するテスト(t *testing.T) {
	type args struct {
		posts []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "IDが１のデータを取得できることを確認すること",
			args: args{
				posts: []interface{}{
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
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startMockDB(tt.args.posts, func(s *Service) {
				var get model.Post
				s.FindOne(&get, 1)
				if get.ID != tt.want {
					t.Errorf("指定したデータが取得できていません get = %v", get)
				}
			})
		})
	}
}
func Testデータベースに保存されている全ての投稿データを取得出来ることを確認するテスト(t *testing.T) {
	type args struct {
		posts []interface{}
		where []interface{}
	}
	type wants struct {
		posts []model.Post
	}
	tests := []struct {
		name    string
		args    args
		wants   wants
		wantErr bool
	}{
		{
			name: "全ての投稿データを取得できていることを確認すること",
			args: args{
				posts: []interface{}{
					&model.Post{
						Name: "Google",
						URL:  "https://www.google.com",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
					&model.Post{
						Name: "Yahoo",
						URL:  "https://www.yahoo.ne.jp",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
					&model.Post{
						Name: "Go",
						URL:  "https://www.golang.org",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
				},
				where: []interface{}{},
			},
			wants: wants{
				posts: []model.Post{
					model.Post{
						Name: "Google",
						URL:  "https://www.google.com",
						Model: model.Model{
							ID:        1,
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
					model.Post{
						Name: "Yahoo",
						URL:  "https://www.yahoo.ne.jp",
						Model: model.Model{
							ID:        2,
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
					model.Post{
						Name: "Go",
						URL:  "https://www.golang.org",
						Model: model.Model{
							ID:        3,
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startMockDB(tt.args.posts, func(s *Service) {
				posts := []model.Post{}
				s.FindAll(&posts, "created_at desc", tt.args.where...)
				jget, _ := json.Marshal(posts)
				jwant, _ := json.Marshal(tt.wants.posts)
				get := string(jget)
				want := string(jwant)
				if get != want {
					t.Errorf("FindAll() get = %v want = %v", get, want)
				}
			})
		})
	}
}

func Testデータベースにデータを保存できることを確認する(t *testing.T) {
	type args struct {
		posts []interface{}
	}
	type wants struct {
		post model.Post
	}
	tests := []struct {
		name    string
		args    args
		wants   wants
		wantErr bool
	}{
		{
			name: "データベースに指定した形式でデータを登録できることを確認すること",
			args: args{
				posts: []interface{}{
					model.Post{
						Name: "Google",
						URL:  "https://www.google.com",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
				},
			},
			wants: wants{
				model.Post{
					Name: "Google",
					URL:  "https://www.google.com",
					Model: model.Model{
						ID:        1,
						CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := database.NewTestHandler()
			defer os.Remove("test.sqlite3")
			defer handler.DB.Close()
			migrate.Set(handler.DB)

			s := Service{
				Handler: handler,
			}

			s.Create(&tt.args.posts)

			var get model.Post
			s.FindOne(&get, 1)

			if reflect.DeepEqual(get, tt.wants.post) {
				t.Errorf("データが保存されていません get = %v want = %v", get, tt.wants.post)
			}
		})
	}
}

func Testデータベースのデータを更新できることを確認するテスト(t *testing.T) {
	type args struct {
		post   model.Post
		params map[string]interface{}
	}
	type wants struct {
		post model.Post
	}
	tests := []struct {
		name    string
		args    args
		wants   wants
		wantErr error
	}{
		{
			name: "Post.Nameに「Google」、Post.URLに「https://www.google.com」と変更を加えたときに、データが更新されることを確認する",
			args: args{
				post: model.Post{Model: model.Model{ID: 1}},
				params: map[string]interface{}{
					"name": "Google",
					"url":  "https://www.google.com",
				},
			},
			wants: wants{
				model.Post{
					Name: "Google",
					URL:  "https://www.google.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPost := []interface{}{
				&model.Post{
					Name: "モックです",
					URL:  "モック用のモデルです",
				},
			}
			startMockDB(mockPost, func(s *Service) {
				s.First(&tt.args.post)
				s.Update(&tt.args.post, tt.args.params)

				var get model.Post
				s.FindOne(&get, 1)
				if get.Name != tt.wants.post.Name || get.URL != tt.wants.post.URL {
					t.Errorf("値が更新できていません get = %v want = %v", get, tt.wants.post)
				}
			})
		})
	}
}

func Testデータを削除できることを確認するテスト(t *testing.T) {
	type args struct {
		posts []interface{}
	}
	type wants struct {
		post model.Post
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "PostモデルのIDが1のデータを削除できることを確認すること",
			args: args{
				posts: []interface{}{
					&model.Post{
						Name: "Google",
						URL:  "https://www.google.com",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startMockDB(tt.args.posts, func(s *Service) {
				post := model.Post{}
				s.Delete(&post, 1)
				var get model.Post
				s.FindOne(&get, 1)

				if get.Name != "" || get.URL != "" {
					t.Errorf("値が削除されていません get = %v", get)
				}
			})
		})
	}
}

func Test選択したカラムのデータのみを取得できることを確認するテスト(t *testing.T) {
	type args struct {
		posts []interface{}
		field string
		where []interface{}
	}
	type wants struct {
		posts []model.Post
	}
	tests := []struct {
		name    string
		args    args
		wants   wants
		wantErr bool
	}{
		{
			name: "PostモデルのNameカラムのみを取得できることを確認する",
			args: args{
				posts: []interface{}{
					&model.Post{
						Name: "Google",
						URL:  "https://www.google.com",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
					&model.Post{
						Name: "Yahoo",
						URL:  "https://www.yahoo.ne.jp",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
					&model.Post{
						Name: "Go",
						URL:  "https://www.golang.org",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
				},
				field: "name",
				where: []interface{}{},
			},
			wants: wants{
				posts: []model.Post{
					model.Post{
						Name: "Google",
					},
					model.Post{
						Name: "Yahoo",
					},
					model.Post{
						Name: "Go",
					},
				},
			},
		},
		{
			name: "PostモデルのURLカラムのみを取得できることを確認する",
			args: args{
				posts: []interface{}{
					&model.Post{
						Name: "Google",
						URL:  "https://www.google.com",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
					&model.Post{
						Name: "Yahoo",
						URL:  "https://www.yahoo.ne.jp",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
					&model.Post{
						Name: "Go",
						URL:  "https://www.golang.org",
						Model: model.Model{
							CreatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
							UpdatedAt: time.Date(2014, 12, 31, 8, 4, 18, 0, loc),
						},
					},
				},
				field: "url",
				where: []interface{}{},
			},
			wants: wants{
				posts: []model.Post{
					model.Post{
						URL: "https://www.google.com",
					},
					model.Post{
						URL: "https://www.yahoo.ne.jp",
					},
					model.Post{
						URL: "https://www.golang.org",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startMockDB(tt.args.posts, func(s *Service) {
				posts := []model.Post{}
				s.Select(tt.args.field).FindAll(&posts, "created_at desc")
				for i := 0; i < len(posts); i++ {
					if !reflect.DeepEqual(posts[i], tt.wants.posts[i]) {
						t.Errorf("Select.FindAll() get = %v want = %v", posts[i], tt.wants.posts[i])
					}
				}
			})
		})
	}
}
