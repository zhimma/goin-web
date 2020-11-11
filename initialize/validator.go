package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	globalInstance "github.com/zhimma/goin-web/global"
	"reflect"
	"strings"
)

// 定义一个全局翻译器T

func Validator(locale string) {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("zh"), ",", 2)[0]
			// fmt.Println(fld.Tag.Get("json"), fld.Tag.Get("zh"))
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		globalInstance.Translator, ok = uni.GetTranslator(locale)
		if !ok {
			globalInstance.SystemLog.Info("设置locale失败")
			return
		}

		// 注册翻译器
		var err error
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, globalInstance.Translator)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, globalInstance.Translator)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, globalInstance.Translator)
		}
		if err != nil {
			globalInstance.SystemLog.Info("设置locale失败")
			return
		}
	}
}
