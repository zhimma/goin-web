package globalInstance

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/zhimma/goin-web/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Viper        *viper.Viper
	BaseConfig   config.Base
	SystemLog    *zap.Logger
	DB           *gorm.DB
	SystemConfig config.System
	Validator    *validator.Validate
)
