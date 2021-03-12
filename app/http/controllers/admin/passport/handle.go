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
	jwtLibrary "github.com/zhimma/goin-web/library/jwt"
	"go.uber.org/zap"
	"strconv"
)

type ClientAuth struct {
	ClientId     string ``
	ClientSecret string ``
}
type ApplyParams struct {
	ContactPhone   string `json:"contact_phone" binding:"required" zh:"联系电话"`
	ContactAddress string `json:"contact_address" binding:"required" zh:"联系地址"`
}

type ClientInfo struct {
	Client model.Client
}

func Apply(c *gin.Context) {
	applyParams := ApplyParams{}
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
	clientId := fmt.Sprintf("%s", globalInstance.UniqueId.Generate())
	// 获取jwt结构体实例
	j := jwtLibrary.JWT{
		SigningKey:          []byte(globalInstance.BaseConfig.Jwt.JwtSecret),
		AccessTokenExpires:  globalInstance.BaseConfig.Jwt.JwtTtl,
		RefreshTokenExpires: globalInstance.BaseConfig.Jwt.JwtRefreshTtl,
	}
	if tokenId, err := strconv.ParseInt(clientId, 10, 64); err != nil {
		response.FailWithMessage("数据转换失败", c)
		globalInstance.SystemLog.Error("数据转换失败", zap.Any("err", err))
		return
	} else {
		tokenDetail, tokenErr := j.GenerateJwtToken(tokenId)
		if tokenErr != nil {
			response.FailWithMessage("TOKEN生成失败", c)
			globalInstance.SystemLog.Error("TOKEN生成失败", zap.Any("err", tokenErr))
			return
		}
		data := model.Client{
			BaseModel:      model.BaseModel{},
			ClientId:       clientId,
			ClientSecret:   tokenDetail.AccessToken,
			Status:         0,
			ContactPhone:   applyParams.ContactPhone,
			ContactAddress: applyParams.ContactAddress,
		}
		if err := service.InsertOne(&data); err != nil {
			response.FailWithMessage("保存用户信息失败", c)
			globalInstance.SystemLog.Error("保存用户信息失败", zap.Any("err", err))
			return
		}
		response.OkWithData(data, c)
	}
}

func Auth(c *gin.Context) {

}
