package initialize

import (
	"github.com/zhimma/goin-web/database"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		database.ColumnTypeMean{},
		database.Admin{},
	)
	if err != nil {
		globalInstance.SYSTERM_LOG.Error("register table failed", zap.Any("err", err))
		os.Exit(0)
	}
	globalInstance.SYSTERM_LOG.Info("migrate table 成功")
}
