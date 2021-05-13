package managerService

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/component"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/constant"
	"github.com/zhimma/goin-web/helper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

type IndexParams struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size" `
	Account  string `json:"account" form:"account" `
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
}

type BindRoleParams struct {
	ManagerId int64 `json:"manager_id" form:"manager_id" binding:"required" zh:"管理员id"`
	RoleId    int64 `json:"role_id" form:"role_id" binding:"required" zh:"角色id"`
}

type ChangePasswordParams struct {
	ManagerId   int64  `json:"manager_id" form:"manager_id" binding:"required" zh:"管理员id"`
	OldPassword string `json:"old_password" form:"old_password" binding:"required" zh:"原密码"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required" zh:"新密码"`
}

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
	Avatar string `json:"avatar" form:"avatar" binding:"-" zh:"头像"`
	Name   string `json:"name" form:"name" binding:"-" zh:"名称"`
	Phone  string `json:"phone" form:"phone" binding:"-" zh:"手机号码"`
	Email  string `json:"email" form:"email" binding:"-" zh:"邮箱"`
	Status int8   `json:"status" form:"status" binding:"-" zh:"状态"`
}

// 管理员列表分页
func ManagerList(params IndexParams) (result CommonDbService.PageResult) {
	condition := CommonDbService.PageStruct{
		Page:         params.Page,
		PageSize:     params.PageSize,
		MapWhere:     nil,
		LikeMapWhere: nil,
	}
	condition.MapWhere = make(map[string]interface{})
	if len(params.Account) > 0 {
		condition.MapWhere["account"] = params.Account
	}
	if len(params.Phone) > 0 {
		condition.LikeMapWhere["phone"] = strings.ToUpper(params.Phone)
	}
	if len(params.Email) > 0 {
		condition.LikeMapWhere["email"] = params.Email
	}
	data := make([]model.Manager, 0, condition.PageSize)
	result = CommonDbService.Paginate(condition, &data)
	return
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

// 绑定用户角色
func BindUserRole(params BindRoleParams) (err error) {
	// 修改用户角色绑定
	record := model.Manager{}
	err = CommonDbService.DetailById(&record, params.ManagerId)
	if err != nil {
		globalInstance.SystemLog.Error("查询用户角色失败", zap.Any("error", err.Error()))
		return errors.New(err.Error())
	}

	if record.IsSuper == 1 {
		return errors.New("超级管理员无需添加权限")
	}
	data := model.Manager{RoleId: params.RoleId}
	_, err = CommonDbService.UpdateByIdNoCheck(record, data)
	if err != nil {
		globalInstance.SystemLog.Error("管理员绑定角色失败", zap.Any("error", err.Error()))
		return errors.New("绑定角色失败")
	}
	return nil
}

// 修改管理员信息
func UpdateManagerInfo(params UpdateManagerParams, managerId int64) (err error) {
	data := model.Manager{
		Avatar: params.Avatar,
		Name:   params.Name,
		Phone:  params.Phone,
		Email:  params.Email,
		Status: params.Status,
	}
	_, err = CommonDbService.UpdateById(&model.Manager{}, managerId, data)
	if err != nil {
		globalInstance.SystemLog.Error("管理员信息修改失败", zap.Any("error", err.Error()))
		return errors.New("管理员信息修改失败")
	}
	return nil
}

// 检查当前用户更新权限
func CheckCurrentUserUpdateHandle(managerId int64, c *gin.Context) (err error) {
	managerRecord, errs := helper.GetCurrentManagerInfo(c)
	if errs != nil {
		return errors.New(errs.Error())
	}
	// 超管可以修改所有 自己只能改自己
	flag := false
	if managerRecord.IsSuper == 1 {
		flag = true
	} else {
		if managerRecord.ID == managerId {
			flag = true
		}
	}
	if !flag {
		return errors.New("没有权限修改")
	}
	return nil
}

// 检查当前用户删除权限
func CheckCurrentUserDeleteHandle(managerId int64, c *gin.Context) (err error) {
	managerRecord, errs := helper.GetCurrentManagerInfo(c)
	if errs != nil {
		return errors.New(errs.Error())
	}
	where := map[string]interface{}{
		"is_super": 1,
	}
	superRecord := model.Manager{}
	// 可优化从缓存中获取
	if err := CommonDbService.DetailByMapWhere(where, &superRecord); err != nil {
		return errors.New("超管信息不存在,请检查")
	}
	flag := true
	// 如果是超管  判断是在删除超管账号 不允许
	// 如果是超管  判断不是在删除超管账号 允许
	// 如果不是超管  有对应权限则可以删除
	if managerRecord.IsSuper == 1 {
		// 不允许删除管理员
		if managerRecord.ID == managerId {
			flag = false
		}
	} else {
		// 防止其他有权限管理员删除超管
		if managerId == superRecord.ID {
			flag = false
		}
	}

	if !flag {
		return errors.New("没有权限删除")
	}
	return nil
}

// 修改密码
func ChangePassword(params ChangePasswordParams) (err error) {
	// 查询用户信息获取record
	record := model.Manager{}
	err = CommonDbService.DetailById(&record, params.ManagerId)
	if err != nil {
		globalInstance.SystemLog.Error("查询用信息失败", zap.Any("error", err.Error()))
		return errors.New(err.Error())
	}
	// 获取盐值
	salt := record.Salt
	// 检查record的密码字段和old_password
	if helper.CompareHashString(record.Password, []byte(params.OldPassword+salt)) == false {
		return errors.New("原密码不正确")
	}
	// 修改用户密码
	password, _ := helper.GenerateHashString(params.NewPassword, salt)
	data := model.Manager{Password: password}
	_, err = CommonDbService.UpdateByIdNoCheck(record, data)
	if err != nil {
		globalInstance.SystemLog.Error("密码修改失败", zap.Any("error", err.Error()))
		return errors.New("密码修改失败")
	}
	return nil
}

// 删除管理员
func DeleteManager(id int64) (err error) {
	record := model.Manager{}
	record.ID = id
	if err := CommonDbService.DeleteById(&record); err != nil {
		globalInstance.SystemLog.Error("删除管理员失败", zap.Any("error", err.Error()))
		return errors.New("删除管理员失败")
	}
	cache := component.NewRedisCache()
	cacheTokenKey1 := fmt.Sprintf(constant.AdminManagerAccessToken, id)
	cacheTokenKey2 := fmt.Sprintf(constant.AdminManagerInfo, id)
	_ = cache.Delete(cacheTokenKey1)
	_ = cache.Delete(cacheTokenKey2)
	return nil
}

// 管理员详情获取
func ManagerDetail(id int64) (data model.Manager, err error) {
	record := model.Manager{}
	if err := CommonDbService.DetailById(&record, id); err != nil {
		globalInstance.SystemLog.Error("获取管理员信息失败", zap.Any("error", err.Error()))
		return record, errors.New("获取管理员信息失败")
	}
	return record, nil
}
