package api_group

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
)

func Index(c *gin.Context) {
	condition := CommonDbService.PageStruct{
		Page:     1,
		PageSize: 1,
		MapWhere: nil,
	}
	condition.MapWhere = make(map[string]interface{})
	condition.MapWhere["status"] = 0
	var data []model.Admin
	result := CommonDbService.Paginate(condition, &data)
	response.OkWithData(result, c)
	return
}

// 新增api group
func Store(c *gin.Context) {
	var apiData = service.ApiGroupStoreParams{}
	if errs := c.ShouldBindJSON(&apiData); errs != nil {
		err, ok := errs.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(errs.Error(), c)
			return
		}

		errorMessageBag := helper.RemoveTopStruct(err.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}

	if err := service.ApiGroupStore(apiData); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
	return

}
func Update(c *gin.Context) {

}
func Destroy(c *gin.Context) {

}
