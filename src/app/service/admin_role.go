package service

import (
	"errors"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/component"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type RoleStoreParams struct {
	Name   string `json:"name" binding:"required" zh:"角色名称"`
	Status int    `json:"status" binding:"required" zh:"启/禁用状态"`
}

type RoleListParams struct {
	Page     int    `json:"page" form:"page" binding:"-"`
	PageSize int    `json:"page_size" form:"page_size" binding:"-"`
	Name     string `json:"name" form:"name" binding:"-"`
}

type RoleApiRelationData struct {
	RoleId int64   `json:"role_id" form:"role_id" binding:"required"`
	Apis   []int64 `json:"apis" form:"apis" binding:"required"`
}

func RoleStore(params RoleStoreParams) error {
	data := model.Role{Name: params.Name, Status: params.Status}
	if err := CommonDbService.InsertOne(&data); err != nil {
		globalInstance.SystemLog.Error("新增角色失败", zap.Any("err", err))
		return errors.New("新增角色失败")
	}
	return nil
}

// 角色列表
func RoleList(params RoleListParams) (result CommonDbService.PageResult) {
	condition := CommonDbService.PageStruct{
		Page:         params.Page,
		PageSize:     params.PageSize,
		MapWhere:     nil,
		LikeMapWhere: nil,
	}
	condition.MapWhere = make(map[string]interface{})
	condition.LikeMapWhere = make(map[string]interface{})
	if len(params.Name) > 0 {
		condition.LikeMapWhere["name"] = params.Name
	}
	var data = make([]model.Role, 0, condition.PageSize)
	result = CommonDbService.Paginate(condition, &data)
	return
}

// 修改保存角色
func RoleUpdate(params RoleStoreParams, id int64) (row int64, err error) {
	modelData := model.Role{
		Name:   params.Name,
		Status: params.Status,
	}
	originModel := model.Role{}
	row, err = CommonDbService.UpdateById(&originModel, id, &modelData)
	if err != nil {
		globalInstance.SystemLog.Error("修改角色失败", zap.Any("err", err))
		return 0, errors.New("修改角色失败")
	}
	return row, nil
}

// 保存角色的api权限
func StoreRoleApis(params RoleApiRelationData) error {
	roleInfo := model.Role{}
	if err := CommonDbService.DetailById(&roleInfo, params.RoleId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在或被删除")
		}
		return errors.New("保存角色的api权限失败【role】")
	}
	apiInfo := make([]model.Api, 0, len(params.Apis))
	if err := CommonDbService.ListByIdSlice(params.Apis, &apiInfo); err != nil {
		return errors.New("保存角色的api权限失败【api】")
	}
	if len(apiInfo) <= 0 {
		return errors.New("api不存在或被删除")
	}
	_, err := addPoliciesWithApis(roleInfo, apiInfo)
	if err != nil {
		return errors.New("保存角色的api权限失败")
	}
	return nil
}

// 给角色添加权限
func addPoliciesWithApis(roleInfo model.Role, apiInfo []model.Api) (result bool, err error) {
	var policies = make([][]string, 0, len(apiInfo))
	for _, v := range apiInfo {
		policies = append(policies, []string{strconv.FormatInt(roleInfo.ID, 10), v.Path, v.Method})
	}
	c := component.Casbin()
	result, err = c.AddPolicies(policies)
	return
}
