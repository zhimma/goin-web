package passport

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service/managerService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
)

// 管理员登陆
func Login(c *gin.Context) {
	loginParams := managerService.LoginParams{}
	if errs := c.ShouldBind(&loginParams); errs != nil {
		err, ok := errs.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(errs.Error(), c)
		}
		errorMessage := helper.RemoveTopStruct(err.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessage[0], c)
		return
	}
	// 去登陆
	adminData := model.Manager{}
	if err := managerService.Login(loginParams, &adminData); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 更新登陆时间 登陆次数 客户端ip
	managerService.UpdateLoginTime(&adminData, c)
	// 缓存登陆token 和 用户基本信息
	if err := managerService.AdminGrantTokenAndCache(&adminData); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 获取缓存信息，准备返回给客户端
	if data, err := managerService.GetManagerInfoFromCache(adminData.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithData(data, c)
		return
	}
}

// 注册用户
func Register(c *gin.Context) {
	// 数据检验
	var u managerService.RegisterData
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
		response.ValidateFail(errorMessageBag[0], c)
		return
	}

	// 判断账号是否已经存在
	checkMap := map[string]interface{}{
		"Account": u.Account,
	}
	if err := managerService.CheckAdminField(checkMap); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := managerService.RegisterUser(u, c); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户注册成功", c)
	return
}
