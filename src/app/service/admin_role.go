package service

import (
	"errors"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
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
