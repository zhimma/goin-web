package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service/adminAuthService"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
)

// 注册用户
func Register(c *gin.Context) {
	// 数据检验
	var u adminAuthService.RegisterData
	if err := c.ShouldBind(&u); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			response.FailWithMessage(err.Error(), c)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.FailWithMessage(errorMessageBag[0], c)
		return
	}

	// 判断账号是否已经存在
	checkMap := map[string]interface{}{
		"Account": u.Account,
	}
	if err := adminAuthService.CheckAdminField(checkMap); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := adminAuthService.RegisterUser(u, c); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户注册成功", c)
	return
}
