package initialize

import (
	"github.com/zhimma/goin-web/database"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		database.ColumnTypeMean{},
		model.Admin{},
		model.Article{},
		model.Category{},
	)
	if err != nil {
		globalInstance.SystemLog.Error("register table failed", zap.Any("err", err))
		os.Exit(0)
	}
	globalInstance.SystemLog.Info("migrate table 成功")
}
