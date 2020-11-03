package globalInstance

import (
	"github.com/spf13/viper"
	"github.com/zhimma/goin-web/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	VIP         *viper.Viper
	BASE_CONFIG config.Base
	SYSTERM_LOG *zap.Logger
	DB          *gorm.DB
	SYSTEM      config.System
)
