package component

import (
	"github.com/zhimma/goin-web/database/seed"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) {
	seed.Manager(db)
}
