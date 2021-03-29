package component

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
)

func Casbin() *casbin.Enforcer {
	mysqlConfig := globalInstance.BaseConfig.Mysql
	adapter, err := gormadapter.NewAdapter(globalInstance.BaseConfig.System.DbType, mysqlConfig.Username+":"+mysqlConfig.Password+"@tcp("+mysqlConfig.Host+")/"+mysqlConfig.Dbname, true)
	if err != nil {
		globalInstance.SystemLog.Error("初始化cabin+gorm驱动失败", zap.Any("err", err))
	}
	enforcer, err := casbin.NewEnforcer(globalInstance.BaseConfig.Casbin.ModelPath, adapter)
	if err != nil {
		globalInstance.SystemLog.Error("初始化cabin NewEnforcer失败", zap.Any("err", err))
	}
	// e.AddFunction("ParamsMatch", ParamsMatchFunc)
	_ = enforcer.LoadPolicy()
	return enforcer
}
