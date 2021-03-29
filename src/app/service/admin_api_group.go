package service

import (
	"errors"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
)

type ApiGroupListParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
}

type ApiGroupStoreParams struct {
	Name        string `json:"name" binding:"required" zh:"组名称"`
	Description string `json:"description" binding:"required" zh:"组描述"`
	Status      int    `json:"status" binding:"required" zh:"组状态"`
}

func ApiGroupStore(apiGroupData ApiGroupStoreParams) error {
	modelData := model.ApiGroup{
		Name:        apiGroupData.Name,
		Description: apiGroupData.Description,
		Status:      apiGroupData.Status,
	}
	if err := CommonDbService.InsertOne(&modelData); err != nil {
		globalInstance.SystemLog.Error("新增apiGroup失败", zap.Any("err", err))
		return errors.New("新增apiGroup失败")
	}
	return nil
}
