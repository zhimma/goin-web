package managers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service/managerService"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
	"net/http"
)

// 列表
func Index(c *gin.Context) {
}

// 新增api
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

	managerId, exist := c.Get("managerId")
	if !exist {
		response.Abort(http.StatusInternalServerError, "用户登陆状态获取失败或获取UID出错", c)
		return
	}
	int64ManagerId, ok := managerId.(int64)
	if !ok {
		response.Abort(http.StatusInternalServerError, "用户登陆状态获取失败或获取UID出错", c)
		return
	}

	if err := managerService.CreateManager(params, int64ManagerId); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("添加管理员成功", c)
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
