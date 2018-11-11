package service

import (
	"errors"

	"github.com/keiya01/eserver/model"
	"github.com/keiya01/eserver/service/database"
)

// PostService DBを扱う構造体で依存関係を分離するために構造体に格納する。
type PostService struct {
	handler *database.Handler
}

// NewPostService PostServiceからDB操作を行うための構造体を作成する。
func NewPostService() *PostService {
	db := database.NewHandler()
	return &PostService{
		handler: db,
	}
}

// FindAll DBからPostに関する全てのデータを取得する。
func (p *PostService) FindAll() (*[]model.Post, error) {
	var posts []model.Post
	p.handler.DB.Order("created_at desc").Find(&posts)
	if len(posts) == 0 {
		err := errors.New("ユーザーが存在していません")
		return &posts, err
	}
	return &posts, nil
}
