package globalInstance

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/zhimma/goin-web/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	BaseConfig config.Base
	SystemLog  *zap.Logger
	DB         *gorm.DB
	Translator ut.Translator
)
