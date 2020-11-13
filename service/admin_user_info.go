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

func AdminUserInfo(id uint) {
	fmt.Println(id)
}

func AdminLogin(admin *model.Admin) (user model.Admin, err error) {
	var data model.Admin
	err = globalInstance.DB.Where("account = ?", admin.Account).Find(&data).Error
	return data, err
}

func CacheAdminUserToken(uid uint, tokenDetail *structure.JwtTokenDetails) error {
	at := time.Unix(tokenDetail.AccessTokenExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(tokenDetail.RefreshTokenExpires, 0)
	now := time.Now()
	fmt.Println(at, rt, now, at.Sub(now))

	accessKey := fmt.Sprintf(constant.AdminUserAccessToken, tokenDetail.AccessTokenUuid)
	errAccess := globalInstance.RedisClient.Set(accessKey, uid, at.Sub(now)).Err()
	if errAccess != nil {
		globalInstance.SystemLog.Error("缓存accessToken失败", zap.Any("error", errAccess))
		return errAccess
	}
	refreshKey := fmt.Sprintf(constant.AdminUserRefreshToken, tokenDetail.RefreshTokenUuid)
	errRefresh := globalInstance.RedisClient.Set(refreshKey, uid, rt.Sub(now)).Err()
	if errRefresh != nil {
		globalInstance.SystemLog.Error("缓存refreshToken失败", zap.Any("error", errRefresh))
		return errRefresh
	}
	return nil
}
