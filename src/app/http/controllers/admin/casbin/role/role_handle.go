package role

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
	"strconv"
)

// 角色列表
func Index(c *gin.Context) {
	params := service.RoleListParams{}
	fmt.Println(c.GetHeader("Content-Type"))
	if errs := c.Bind(&params); errs != nil {
		validateErr, ok := errs.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(errs.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(validateErr.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	fmt.Println(params)
	data := service.RoleList(params)
	response.OkWithData(data, c)
	return

}

// 新增角色
func Store(c *gin.Context) {
	params := service.RoleStoreParams{}
	if errs := c.ShouldBind(&params); errs != nil {
		err, ok := errs.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(errs.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(err.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	if err := service.RoleStore(params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("新增成功", c)
	return
}

func Update(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		response.ValidateFail("获取参数id失败", c)
		var roleStoreParams = service.RoleStoreParams{}
		if errs := c.ShouldBind(&roleStoreParams); errs != nil {
			err, ok := errs.(validator.ValidationErrors)
			if !ok {
				response.FailWithMessage(errs.Error(), c)
				return
			}
			errorMessageBag := helper.RemoveTopStruct(err.Translate(globalInstance.Translator))
			response.ValidateFail(errorMessageBag[0], c)
			return
		}
		_, err = service.RoleUpdate(roleStoreParams, id)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.OkWithMessage("修改成功", c)
		return
	}
}

func Destroy(c *gin.Context) {

}
