package role

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
)

func Index(c *gin.Context) {

}

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

}

func Destroy(c *gin.Context) {

}
