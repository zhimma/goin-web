package main

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/zhimma/goin-web/core"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/initialize"
)

var trans ut.Translator

func main() {
	// 加载配置文件
	globalInstance.Viper = core.Viper()
	// 注册日志系统
	globalInstance.SystemLog = core.Zap()

	// 加载数据验证器
	initialize.Validator("zh")
	// 注册mysql
	globalInstance.DB = initialize.Gorm()
	db, _ := globalInstance.DB.DB()

	// 执行数据库迁移
	initialize.Migrate(globalInstance.DB)
	// 初始化执行sql seed
	// initialize.Seeder(globalInstance.DB)
	defer db.Close()
	// 	启动服务
	core.Run()
}
