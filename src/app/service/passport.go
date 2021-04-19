package service

import (
	"encoding/json"
	"fmt"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	"github.com/zhimma/goin-web/database/structure"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/constant"
	"github.com/zhimma/goin-web/helper"
	jwtLibrary "github.com/zhimma/goin-web/library/jwt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type AuthParams struct {
	ClientId     int64  `json:"client_id" binding:"required" zh:"client_id"`
	ClientSecret string `json:"client_secret" binding:"required" zh:"client_secret"`
}
type ApplyParams struct {
	ContactPhone   string `json:"contact_phone" binding:"required" zh:"联系电话"`
	ContactAddress string `json:"contact_address" binding:"required" zh:"联系地址"`
}

// 申请登陆
func ApplyClient(params ApplyParams) (data model.Client, error error) {
	clientId := globalInstance.UniqueId.Generate().Int64()
	salt := helper.RandStringBytes(4)
	clientSecret, _ := helper.GenerateHashString(strconv.FormatInt(clientId, 10), salt)
	data = model.Client{
		BaseModel:      model.BaseModel{},
		ClientId:       clientId,
		Salt:           salt,
		ClientSecret:   clientSecret,
		Status:         0,
		ContactPhone:   params.ContactPhone,
		ContactAddress: params.ContactAddress,
	}
	where := map[string]interface{}{
		"contact_phone":   params.ContactPhone,
		"contact_address": params.ContactAddress,
	}
	if err := CommonDbService.FirstOrCreate(where, &data); err != nil {
		globalInstance.SystemLog.Error("保存用户信息失败", zap.Any("err", err))
		return model.Client{}, err
	}
	return data, nil
}

// 缓存客户端信息
func CacheClientInfo(client *model.Client) error {
	config := globalInstance.BaseConfig.Jwt
	at := time.Unix(config.JwtTtl, 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	tokenData, makeTokenErr := makeToken(client.ClientId)
	if makeTokenErr != nil {
		globalInstance.SystemLog.Error("生成token失败", zap.Any("error", makeTokenErr))
		return makeTokenErr
	}
	cacheTokenKey := fmt.Sprintf(constant.ClientAuthToken, client.ClientId)
	cacheTokenError := globalInstance.RedisClient.Set(cacheTokenKey, tokenData.Token, at.Sub(now)).Err()
	if cacheTokenError != nil {
		globalInstance.SystemLog.Error("缓存accessToken失败", zap.Any("error", cacheTokenError))
		return cacheTokenError
	}
	jsonData, err := json.Marshal(client)
	if err != nil {
		globalInstance.SystemLog.Error("缓存clientInfo失败「to json」", zap.Any("error", cacheTokenError))
		return err
	}
	cacheInfoKey := fmt.Sprintf(constant.ClientInfo, client.ClientId)
	cacheClientInfo := globalInstance.RedisClient.Set(cacheInfoKey, jsonData, at.Sub(now)).Err()
	if cacheClientInfo != nil {
		globalInstance.SystemLog.Error("缓存clientInfo失败", zap.Any("error", cacheClientInfo))
		return cacheClientInfo
	}
	return nil
}

// 生成客户端登陆token
func makeToken(identifier int64) (tokenData *structure.JwtTokenDetails, err error) {
	jwt := jwtLibrary.NewJWT()
	return jwt.GenerateJwtToken(identifier)
}

// 从缓存中获取客户端信息
func GetClientInfoFromCache(clientId int64) (data map[string]interface{}, err error) {
	var clientInfo model.Client
	cacheTokenKey := fmt.Sprintf(constant.ClientAuthToken, clientId)
	tokenData, getTokenCacheError := globalInstance.RedisClient.Get(cacheTokenKey).Result()
	if getTokenCacheError != nil {
		globalInstance.SystemLog.Error("获取accessToken缓存失败", zap.Any("error", getTokenCacheError))
		return data, getTokenCacheError
	}
	cacheInfoKey := fmt.Sprintf(constant.ClientInfo, clientId)
	clientData, cacheClientInfo := globalInstance.RedisClient.Get(cacheInfoKey).Result()
	if cacheClientInfo != nil {
		globalInstance.SystemLog.Error("获取clientInfo缓存失败", zap.Any("error", cacheClientInfo))
		return data, getTokenCacheError
	}
	var resultData = make(map[string]interface{})
	if err := json.Unmarshal([]byte(clientData), &clientInfo); err != nil {
		return resultData, err
	}
	resultData["tokenInfo"] = tokenData
	resultData["clientInfo"] = clientInfo
	return resultData, nil
}
