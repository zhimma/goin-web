package client_passport

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
)

// 申请入驻
func Apply(c *gin.Context) {
	applyParams := service.ApplyParams{}
	if err := c.ShouldBindJSON(&applyParams); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ValidateFail(err.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	if data, err := service.ApplyClient(applyParams); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}

}
