package passport

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// 申请入驻
func Apply(c *gin.Context) {
	applyParams := service.ApplyParams{}
	if err := c.ShouldBindJSON(&applyParams); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(err.Error(), c)
			return
		}
		fmt.Println(errs.Translate(globalInstance.Translator))
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.FailWithMessage(errorMessageBag[0], c)
		return
	}
	if data, err := service.ApplyClient(applyParams); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

// auth 登陆
func Auth(c *gin.Context) {
	authData := service.AuthParams{}
	if err := c.ShouldBindJSON(&authData); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(errs.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.FailWithMessage(errorMessageBag[0], c)
		return
	}
	where := make(map[string]interface{})
	where["client_id"] = authData.ClientId
	where["client_secret"] = authData.ClientSecret
	modelRecord := model.Client{}
	if err := service.DetailByWhere(where, &modelRecord); err != nil {
		globalInstance.SystemLog.Error("查询client_id信息失败", zap.Error(err))
		return
	}
	if modelRecord.Status == 0 {
		response.Abort(http.StatusForbidden, "Unauthorized:Disabled", c)
		return
	}
	if modelRecord.Status == 2 {
		response.Abort(http.StatusForbidden, "Unauthorized:Forbidden", c)
		return
	}
	if err := service.CacheClientInfo(&modelRecord); err != nil {
		globalInstance.SystemLog.Error("生成client缓存失败", zap.Error(err))
		return
	}
	if data, err := service.GetClientInfoFromCache(modelRecord.ClientId); err != nil {
		globalInstance.SystemLog.Error("读取client缓存失败", zap.Error(err))
		return
	} else {
		ttl := globalInstance.BaseConfig.Jwt.JwtTtl
		result := map[string]interface{}{
			"token":          data["tokenInfo"],
			"expires_time":   time.Now().Add(time.Duration(ttl)),
			"expires_second": ttl,
		}
		response.OkWithData(result, c)
		return
	}
}
