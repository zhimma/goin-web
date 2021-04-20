package managers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service/managerService"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
)

// 列表
func Index(c *gin.Context) {
	var indexData managerService.IndexParams
	if errs := c.ShouldBind(&indexData); errs != nil {
		err, ok := errs.(validator.ValidationErrors)
		if !ok {
			response.ValidateFail(errs.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(err.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	result := managerService.ManagerList(indexData)
	response.OkWithData(result, c)
	return
}

// 新增管理员
func Store(c *gin.Context) {
	var params managerService.CreateManagerParams
	if err := c.ShouldBind(&params); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(err.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}

	managerRecord, err := helper.GetCurrentManagerInfo(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := managerService.CreateManager(params, managerRecord.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("添加管理员成功", c)
	return
}

func BindRole(c *gin.Context) {
	var params managerService.BindRoleParams
	if bindError := c.ShouldBind(&params); bindError != nil {
		validatorErr, ok := bindError.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(bindError.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(validatorErr.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	err := managerService.BindUserRole(params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户角色绑定成功", c)
	return
}

// 查询
func Show(c *gin.Context) {

}

// 更新
func Update(c *gin.Context) {
}

// 删除
func Destroy(c *gin.Context) {
}

// 修改密码
func ChangePassword(c *gin.Context) {

	var params managerService.ChangePasswordParams
	if bindError := c.ShouldBind(&params); bindError != nil {
		validatorError, ok := bindError.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(bindError.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(validatorError.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	managerRecord, err := helper.GetCurrentManagerInfo(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 超管可以修改所有密码 自己只能改自己的密码
	flag := false
	if managerRecord.IsSuper == 1 {
		flag = true
	} else {
		if managerRecord.ID == params.ManagerId {
			flag = true
		}
	}
	if !flag {
		response.FailWithMessage("没有权限修改", c)
		return
	}
	if err := managerService.ChangePassword(params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("密码修改成功", c)
	return
}
