package adminAuthService

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/helper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type RegisterData struct {
	Account    string `json:"account" form:"account" binding:"required" zh:"账号"`
	Password   string `json:"password" form:"account" binding:"required" zh:"密码"`
	RePassword string `json:"RePassword" form:"RePassword" binding:"required,eqfield=Password" zh:"重复密码"`
}

// 检查字段是否存在
func CheckAdminField(where map[string]interface{}) error {
	var data model.Admin
	err := CommonDbService.DetailByMapWhere(where, &data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return errors.New("用户已存在")
	}
}

// 注册用户service
func RegisterUser(params RegisterData, c *gin.Context) error {
	// 注册用户
	salt := helper.RandStringBytes(4)
	password, _ := helper.GenerateHashString(params.Password, salt)
	data := model.Admin{
		Account:     params.Account,
		Salt:        salt,
		Password:    password,
		Avatar:      "",
		Name:        params.Account,
		Phone:       params.Account,
		Status:      0,
		LastLoginIp: c.ClientIP(),
		LastLoginAt: time.Now(),
		LoginTimes:  0,
	}
	if err := CommonDbService.InsertOne(&data); err != nil {
		globalInstance.SystemLog.Error("注册用户写入数据库失败", zap.Any("error", err))
		return errors.New("注册用户失败「to save」")
	}
	return nil
}
