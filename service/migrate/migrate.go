package migrate

import (
	"github.com/jinzhu/gorm"
	"github.com/keiya01/eserver/model"
)

func Set(db *gorm.DB) {
	db.AutoMigrate(&model.Post{})
}
