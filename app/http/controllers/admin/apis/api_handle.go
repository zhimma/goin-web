package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
	"strconv"
)

func Index(c *gin.Context) {
	var indexData = service.ApiListParams{}
	if errs := c.ShouldBindBodyWith(&indexData, binding.JSON); errs != nil {
		err, ok := errs.(validator.ValidationErrors)
		if !ok {
			response.ValidateFail(err.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(err.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	result := service.ApiList(indexData)
	response.OkWithData(result, c)
	return
}

// 新增api
func Store(c *gin.Context) {
	var apiData = service.ApiParams{}
	if errs := c.ShouldBindJSON(&apiData); errs != nil {
		err, ok := errs.(validator.ValidationErrors)
		if !ok {
			response.ValidateFail(err.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(err.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	if err := service.ApiStore(apiData); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
	return
}

// 更新
func Update(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		response.ValidateFail("获取参数id失败", c)
	}
	var apiData = service.ApiParams{}
	if errs := c.ShouldBindJSON(&apiData); errs != nil {
		err, ok := errs.(validator.ValidationErrors)
		if !ok {
			response.ValidateFail(err.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(err.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	_, err = service.ApiUpdate(apiData, id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
	return

}
func Destroy(c *gin.Context) {

}
