package client_passport

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// auth 登陆
func Auth(c *gin.Context) {
	authData := service.AuthParams{}
	if err := c.ShouldBind(&authData); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ValidateFail(errs.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	where := make(map[string]interface{})
	where["client_id"] = authData.ClientId
	where["client_secret"] = authData.ClientSecret
	modelRecord := model.Client{}
	if err := CommonDbService.DetailByMapWhere(where, &modelRecord); err != nil {
		globalInstance.SystemLog.Error("查询client_id信息失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	if data, err := service.GetClientInfoFromCache(modelRecord.ClientId); err != nil {
		globalInstance.SystemLog.Error("读取client缓存失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
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
