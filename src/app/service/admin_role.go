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

func RoleStore(params RoleStoreParams) error {
	data := model.Role{Name: params.Name, Status: params.Status}
	if err := CommonDbService.InsertOne(&data); err != nil {
		globalInstance.SystemLog.Error("新增角色失败", zap.Any("err", err))
		return errors.New("新增角色失败")
	}
	return nil
}
