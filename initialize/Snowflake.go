package initialize

import (
	"github.com/bwmarrin/snowflake"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
)

func SnowFlake() {
	instance, err := snowflake.NewNode(1)
	if err != nil {
		globalInstance.SystemLog.Error("snowflake server register failed, err:", zap.Any("err", err))
	} else {
		globalInstance.UniqueId = instance
	}
}
