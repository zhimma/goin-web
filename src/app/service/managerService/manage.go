package managerService

import (
	"errors"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/helper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type CreateManagerParams struct {
	Account  string `json:"account" form:"account" binding:"required" zh:"账号"`
	Password string `json:"password" form:"password" binding:"required" zh:"密码"`
	Avatar   string `json:"avatar" form:"avatar" binding:"required" zh:"头像"`
	Name     string `json:"name" form:"name" binding:"required" zh:"名称"`
	Phone    string `json:"phone" form:"phone" binding:"required" zh:"手机号码"`
	Email    string `json:"email" form:"email" binding:"required" zh:"邮箱"`
	Status   int8   `json:"status" form:"status" binding:"required" zh:"状态"`
}

type UpdateManagerParams struct {
	Password string `json:"password" form:"password" binding:"required" zh:"密码"`
	Avatar   string `json:"avatar" form:"avatar" binding:"required" zh:"头像"`
	Name     string `json:"name" form:"name" binding:"required" zh:"名称"`
	Phone    string `json:"phone" form:"phone" binding:"required" zh:"手机号码"`
	Email    string `json:"email" form:"email" binding:"required" zh:"邮箱"`
	Status   int8   `json:"status" form:"status" binding:"required" zh:"状态"`
}

// 新增管理员
func CreateManager(params CreateManagerParams, operatorId int64) error {
	maps := []map[string]interface{}{
		{"account": params.Account},
		{"email": params.Email},
		{"phone": params.Phone},
	}
	modelData := model.Manager{}
	// 如果正常返回 则认为用户邮箱、手机号、或者账号已经被使用，查询到了数据
	err := CommonDbService.DetailByMapOrWhere(maps, &modelData)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("新增用户失败,邮箱、手机号、或者账号已经被使用")
	}
	salt := helper.RandStringBytes(4)
	password, _ := helper.GenerateHashString(params.Password, salt)
	data := model.Manager{
		Account:     params.Account,
		Email:       params.Email,
		Salt:        salt,
		Password:    password,
		Avatar:      params.Avatar,
		Name:        params.Name,
		Phone:       params.Phone,
		Status:      params.Status,
		LastLoginAt: time.Now(),
		LoginTimes:  0,
		OperatorId:  operatorId,
	}
	if err := CommonDbService.InsertOne(&data); err != nil {
		globalInstance.SystemLog.Error("新增用户写入数据库失败", zap.Any("error", err))
		return errors.New("新增用户失败「to save」")
	}
	return nil
}