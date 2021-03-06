/**

 * @Author: zhimma
 * @Description:
 * @File:  category
 * @Date: 2020/11/17 16:27
 */
package category

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
	"strconv"
)

type CreateCategory struct {
	Name        string `json:"name" form:"name" binding:"required" zh:"账号"`
	Pid         uint   `json:"pid" form:"pid" binding:"gte=0" zh:"上级分类"`
	Description string `json:"description" form:"description" binding:"required" zh:"分类描述"`
}

type UpdateCategory struct {
	Id          uint   `json:"id" form:"id" binding:"gt=0" zh:"id"`
	Name        string `json:"name" form:"name" binding:"required" zh:"账号"`
	Pid         uint   `json:"pid" form:"pid" binding:"gte=0" zh:"上级分类"`
	Description string `json:"description" form:"description" binding:"required" zh:"分类描述"`
}

func Index(c *gin.Context) {
	where := map[string]interface{}{}
	var categories []model.Category
	if err := CommonDbService.List(where, &categories); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// result := helper.List2Tree(categories, "id", 0, "_child", 0)
	response.OkWithData(categories, c)
}

// 	分类详情
func Show(c *gin.Context) {
	intId := c.Param("id")
	id, _ := strconv.ParseInt(intId, 10, 64)
	data := model.Category{}
	if err := CommonDbService.DetailById(&data, id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
	return
}

// 创建分类
func Store(c *gin.Context) {
	var params CreateCategory
	if err := c.ShouldBindWith(&params, binding.JSON); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(err.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.FailWithMessage(errorMessageBag[0], c)
		return
	}
	data := model.Category{
		Pid:         int64(params.Pid),
		Name:        params.Name,
		Description: params.Description,
	}
	// 公共写入
	if err := CommonDbService.InsertOne(&data); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
	return

}

// 更新分类信息
func Update(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.ParseInt(stringId, 10, 64)
	var params CreateCategory
	// 接受参数
	if err := c.ShouldBind(&params); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ValidateFail(err.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.ValidateFail(errorMessageBag[0], c)
		return
	}
	originModel := model.Category{}
	updateData := map[string]interface{}{
		"pid":         params.Pid,
		"name":        params.Name,
		"description": params.Description,
	}
	num, err := CommonDbService.UpdateById(originModel, id, updateData)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	response.OkWithMessage(fmt.Sprintf("success【%d】", num), c)
}

// 删除
func Destroy(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data := model.Category{
		BaseModel: model.BaseModel{ID: int64(id)},
	}
	if err := CommonDbService.DeleteById(&data); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	response.OkWithMessage("删除成功", c)
}
