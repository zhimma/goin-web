package component

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
		model.Client{},
		model.Api{},
		model.ApiGroup{},
	)

	/*if db.Migrator().HasTable(&model.Client{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci comment '客户端表'").AutoMigrate(&model.Client{})
	}*/

	if err != nil {
		globalInstance.SystemLog.Error("register table failed", zap.Any("err", err))
		os.Exit(0)
	}
	globalInstance.SystemLog.Info("migrate table 成功")
}
