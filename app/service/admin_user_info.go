package service

import (
	"fmt"
	"github.com/zhimma/goin-web/database/constant"
	"github.com/zhimma/goin-web/database/model"
	"github.com/zhimma/goin-web/database/structure"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
	"time"
)

// 检查redis中是否存在用户的token
func AdminUserTokenCheck(tokenInfo *structure.JwtClaims) (string, error) {
	key := fmt.Sprintf(constant.AdminUserAccessToken, tokenInfo.IDENTIFIER, tokenInfo.UUID)
	return globalInstance.RedisClient.Get(key).Result()
}

// 登陆时查询用户信息 by account
func AdminLogin(admin *model.Admin) (user model.Admin, err error) {
	var data model.Admin
	err = globalInstance.DB.Where("account = ?", admin.Account).Find(&data).Error
	return data, err
}

// 注册 创建用户
func AdminRegister(admin *model.Admin) (err error) {
	err = globalInstance.DB.Create(admin).Error
	return
}

// 检查字段是否存在
func CheckAdminField(where map[string]interface{}) (status bool, err error) {
	var data model.Admin
	_ = globalInstance.DB.Where(where).First(&data).Error
	/*if err != nil {
		return false, err
	}*/
	if data.ID != 0 {
		return true, nil
	}
	return false, nil
}

// 缓存token access_token & refresh_token
func CacheAdminUserToken(uid uint, tokenDetail *structure.JwtTokenDetails) error {
	at := time.Unix(tokenDetail.Expires, 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	fmt.Println(at, now, at.Sub(now))

	accessKey := fmt.Sprintf(constant.AdminUserAccessToken, uid, tokenDetail.Uuid)
	errAccess := globalInstance.RedisClient.Set(accessKey, uid, at.Sub(now)).Err()
	if errAccess != nil {
		globalInstance.SystemLog.Error("缓存accessToken失败", zap.Any("error", errAccess))
		return errAccess
	}
	return nil
}
