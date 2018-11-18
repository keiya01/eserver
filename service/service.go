package service

import (
	"github.com/keiya01/eserver/service/database"
	"github.com/pkg/errors"
)

// Service DBを扱う構造体で依存関係を分離するために構造体に格納する。
type Service struct {
	*database.Handler
}

// NewService ServiceからDB操作を行うための構造体を作成する。
func NewService(db *database.Handler) *Service {
	return &Service{db}
}

func (s *Service) Select(field string) *Service {
	s.Handler.DB = s.Handler.Select(field)
	return s
}

func (s *Service) FindOne(model interface{}, where ...interface{}) error {
	if db := s.First(model, where...); db.Error != nil {
		return errors.Wrap(db.Error, "FindOne()")
	}

	return nil
}

// FindAll DBからPostに関する全てのデータを取得する。
func (s *Service) FindAll(model interface{}, order string, where ...interface{}) error {
	if db := s.Order(order).Find(model, where...); db.Error != nil {
		return errors.Wrap(db.Error, "FindAll()")
	}

	return nil
}

// Create 構造体にのフィールドに含むデータをDBに保存する
func (s *Service) Create(model interface{}) error {
	if db := s.Handler.Create(model); db.Error != nil {
		return errors.Wrap(db.Error, "Create()")
	}

	return nil
}

// Update 第一引き数に更新したい構造体を入れて、
// 第二引数にフィールド名と変更する値がセットのマップを入れる
func (s *Service) Update(model interface{}, params map[string]interface{}) error {
	if db := s.Model(model).Updates(params); db.Error != nil {
		return errors.Wrap(db.Error, "Update()")
	}

	return nil
}

// Delete 指定された構造体の情報と一致するデータを削除します。
func (s *Service) Delete(model interface{}) error {
	if db := s.Handler.Delete(model); db.Error != nil {
		return errors.Wrap(db.Error, "Delete()")
	}

	return nil
}
