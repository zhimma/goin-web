package service

import (
	"errors"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type ApiListParams struct {
	Page       int    `json:"page" form:"page" binding:"-"`
	PageSize   int    `json:"page_size" form:"page_size" binding:"-"`
	Path       string `json:"path" form:"path" binding:"-"`
	Method     string `json:"method" form:"method" binding:"-"`
	ApiGroupId int    `json:"api_group_id" form:"api_group_id" binding:"-"`
}

type ApiParams struct {
	Path        string `json:"path" binding:"required" zh:"接口地址"`
	Method      string `json:"method" binding:"required" zh:"请求方法"`
	Description string `json:"description" binding:"required" zh:"接口描述"`
	ApiGroupId  int64  `json:"api_group_id" binding:"required" zh:"api所属组"`
}

func ApiList(params ApiListParams) (result CommonDbService.PageResult) {
	condition := CommonDbService.PageStruct{
		Page:         params.Page,
		PageSize:     params.PageSize,
		MapWhere:     nil,
		LikeMapWhere: nil,
	}
	condition.MapWhere = make(map[string]interface{})
	condition.LikeMapWhere = make(map[string]interface{})

	if params.ApiGroupId > 0 {
		condition.MapWhere["api_group_id"] = params.ApiGroupId
	}
	if len(strings.ToUpper(params.Method)) > 0 {
		condition.LikeMapWhere["method"] = strings.ToUpper(params.Method)
	}
	if len(params.Path) > 0 {
		condition.LikeMapWhere["path"] = params.Path
	}
	var data = make([]model.Api, 0, condition.PageSize)
	result = CommonDbService.Paginate(condition, &data)
	return
}

// 保存api
func ApiStore(params ApiParams) error {
	modelData := model.Api{
		Path:        params.Path,
		Description: params.Description,
		Method:      strings.ToUpper(params.Method),
		ApiGroupId:  params.ApiGroupId,
	}

	isExist, err := apiIsExist(params)
	if err != nil {
		return errors.New(err.Error())
	}
	if isExist {
		return errors.New("api已存在")
	}
	if err := CommonDbService.InsertOne(&modelData); err != nil {
		globalInstance.SystemLog.Error("新增api失败", zap.Any("err", err))
		return errors.New("新增api失败")
	}
	return nil
}

// 保存api
func ApiUpdate(params ApiParams, id int64) (row int64, err error) {
	modelData := model.Api{
		Path:        params.Path,
		Description: params.Description,
		Method:      strings.ToUpper(params.Method),
		ApiGroupId:  params.ApiGroupId,
	}
	originModel := model.Api{}
	row, err = CommonDbService.UpdateById(&originModel, id, modelData)
	if err != nil {
		globalInstance.SystemLog.Error("修改api失败", zap.Any("err", err))
		return 0, errors.New("修改api失败")
	}
	return row, nil
}

// api是否存在
func apiIsExist(params ApiParams) (isExist bool, err error) {
	mapWhere := map[string]interface{}{
		"path":   params.Path,
		"method": strings.ToUpper(params.Method),
	}
	var apiData = model.Api{}
	if err := CommonDbService.DetailByMapWhere(mapWhere, &apiData); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, err
	}

	if apiData.ID >= 0 {
		return true, nil
	}
	return false, err
}
