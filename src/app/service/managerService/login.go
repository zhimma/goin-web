package managerService

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/component"
	"github.com/zhimma/goin-web/database/model"
	"github.com/zhimma/goin-web/database/structure"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/constant"
	"github.com/zhimma/goin-web/helper"
	jwtLibrary "github.com/zhimma/goin-web/library/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type LoginParams struct {
	Account  string `json:"account" form:"account" binding:"required" zh:"账号"`
	Password string `json:"password" form:"password" binding:"required" zh:"密码"`
}

type PassportResult struct {
	ManagerInfo model.Manager          `json:"manager_info"`
	TokenInfo   map[string]interface{} `json:"token_info"`
}

// 登陆认证
func Login(params LoginParams, data *model.Manager) error {
	mapWhere := map[string]interface{}{
		"account": params.Account,
	}
	err := CommonDbService.DetailByMapWhere(mapWhere, data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("账号信息不存在")
	}
	if data.Status == 0 {
		return errors.New("账号已禁用或不存在")
	}
	if helper.CompareHashString(data.Password, []byte(params.Password+data.Salt)) == false {
		return errors.New("密码错误，请检查后重新登陆")
	}
	return nil
}

// 生成token 并缓存
func AdminGrantTokenAndCache(data *model.Manager) error {
	// 获取配置
	cache := component.NewRedisCache()
	config := globalInstance.BaseConfig.Jwt
	// 获取token过期时间
	at := time.Unix(config.JwtTtl, 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	// 生成token
	tokenData, makeTokenErr := makeToken(data.ID)
	if makeTokenErr != nil {
		globalInstance.SystemLog.Error("生成accessToken失败", zap.Any("error", makeTokenErr))
		return errors.New("生成accessToken失败")
	}
	cacheTokenKey := fmt.Sprintf(constant.AdminManagerAccessToken, data.ID)
	cacheTokenError := cache.Set(cacheTokenKey, tokenData.Token, at.Sub(now))
	if cacheTokenError != nil {
		globalInstance.SystemLog.Error("缓存accessToken失败", zap.Any("error", cacheTokenError))
		return errors.New("缓存accessToken失败")

	}
	cacheInfoKey := fmt.Sprintf(constant.AdminManagerInfo, data.ID)
	cacheClientInfoErr := cache.Set(cacheInfoKey, data, at.Sub(now))
	if cacheClientInfoErr != nil {
		globalInstance.SystemLog.Error("缓存用户信息失败", zap.Any("error", cacheClientInfoErr))
		return errors.New("缓存用户信息失败")
	}
	return nil
}

// 生成客户端登陆token
func makeToken(identifier int64) (tokenData *structure.JwtTokenDetails, err error) {
	jwt := jwtLibrary.NewJWT()
	return jwt.GenerateJwtToken(identifier)
}

// 从缓存中获取客户端信息
func GetManagerInfoFromCache(userId int64) (data PassportResult, err error) {
	cache := component.NewRedisCache()

	tokenInfo := structure.JwtTokenDetails{}
	cacheTokenKey := fmt.Sprintf(constant.AdminManagerAccessToken, userId)
	getTokenCacheError := cache.Get(cacheTokenKey, &tokenInfo.Token)
	if getTokenCacheError != nil {
		globalInstance.SystemLog.Error("获取accessToken缓存失败", zap.Any("error", getTokenCacheError))
		return data, errors.New("获取accessToken缓存失败")
	}

	userInfo := model.Manager{}
	cacheInfoKey := fmt.Sprintf(constant.AdminManagerInfo, userId)
	getTokenCacheError = cache.Get(cacheInfoKey, &userInfo)
	if getTokenCacheError != nil {
		globalInstance.SystemLog.Error("获取用户缓存失败", zap.Any("error", getTokenCacheError))
		return data, errors.New("获取用户缓存失败")
	}
	resultData := PassportResult{}
	ttl := globalInstance.BaseConfig.Jwt.JwtTtl
	tokenData := map[string]interface{}{
		"tokenInfo":      tokenInfo.Token,
		"expires_time":   time.Now().Add(time.Duration(ttl)),
		"expires_second": ttl,
	}
	resultData.TokenInfo = tokenData
	resultData.ManagerInfo = userInfo
	return resultData, nil
}

// 更新登陆信息
func UpdateLoginTime(data *model.Manager, c *gin.Context) {
	data.LoginTimes += 1
	data.LastLoginAt = time.Now()
	data.LastLoginIp = c.ClientIP()
	updateData := map[string]interface{}{
		"login_times":   data.LoginTimes,
		"last_login_at": data.LastLoginAt,
		"last_login_ip": data.LastLoginIp,
	}
	if _, err := CommonDbService.UpdateById(data, data.ID, updateData); err != nil {
		globalInstance.SystemLog.Error("更新用户登陆次数相关信息失败", zap.Any("error", err))
	}
}

// 检查redis中是否存在用户的token
func AdminUserTokenCheck(tokenInfo *structure.JwtClaims) (string, error) {
	key := fmt.Sprintf(constant.AdminManagerAccessToken, tokenInfo.IDENTIFIER)
	token := structure.JwtTokenDetails{}
	cache := component.NewRedisCache()
	return token.Token, cache.Get(key, &token.Token)
}
