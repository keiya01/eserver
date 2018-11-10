package migrate

import (
	"github.com/jinzhu/gorm"
	"github.com/keiya01/eserver/model"
)

func Set(g *gorm.DB) {
	g.AutoMigrate(&model.Post{})
}
