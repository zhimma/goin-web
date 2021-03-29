package main

import (
	"github.com/zhimma/goin-web/component"
	"github.com/zhimma/goin-web/core"
	globalInstance "github.com/zhimma/goin-web/global"
)

func main() {
	// 加载配置文件
	core.Viper()
	// 注册日志系统
	globalInstance.SystemLog = core.Zap()
	// 加载数据验证器
	component.Validator("zh")
	// 注册mysql
	globalInstance.DB = component.Gorm()
	db, _ := globalInstance.DB.DB()
	// 执行数据库迁移
	component.Migrate(globalInstance.DB)
	// 初始化执行sql seed
	component.Seeder(globalInstance.DB)
	defer db.Close()

	// 注册redis
	component.RedisClient()

	// 注册雪花算法服务
	component.SnowFlake()

	// 	启动服务
	core.Run()
}
