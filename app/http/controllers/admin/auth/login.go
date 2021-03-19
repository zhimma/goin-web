package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service/adminAuthService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
)

// 管理员登陆
func Login(c *gin.Context) {
	loginParams := adminAuthService.LoginParams{}
	if errs := c.ShouldBindJSON(&loginParams); errs != nil {
		err, ok := errs.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(err.Error(), c)
		}
		errorMessage := helper.RemoveTopStruct(err.Translate(globalInstance.Translator))
		response.FailWithMessage(errorMessage[0], c)
		return
	}
	// 去登陆
	adminData := model.Admin{}
	if err := adminAuthService.Login(loginParams, &adminData); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 更新登陆时间 登陆次数 客户端ip
	adminAuthService.UpdateLoginTime(&adminData, c)
	// 缓存登陆token 和 用户基本信息
	if err := adminAuthService.AdminGrantTokenAndCache(&adminData); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 获取缓存信息，准备返回给客户端
	if data, err := adminAuthService.GetAdminInfoFromCache(adminData.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithData(data, c)
		return
	}
}
