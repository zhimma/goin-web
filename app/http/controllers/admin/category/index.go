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
	"github.com/zhimma/goin-web/app/service"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
)

type CreateCategory struct {
	Name        string `json:"name" form:"name" binding:"required" zh:"账号"`
	Pid         uint   `json:"pid" form:"pid" binding:"gte=0" zh:"上级分类"`
	Description string `json:"description" form:"description" binding:"required" zh:"分类描述"`
}

func Index(c *gin.Context) {
	response.Ok(c)
}

func Show(c *gin.Context) {
}

func Store(c *gin.Context) {
	var params CreateCategory
	if err := c.ShouldBindWith(&params, binding.JSON); err != nil {
		fmt.Println(err.Error(), err)
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.FailWithMessage(err.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.FailWithMessage(errorMessageBag, c)
		return
	}
	data := model.Category{
		Pid:         params.Pid,
		Name:        params.Name,
		Description: params.Description,
	}
	// 公共写入
	if err := service.Create(&data); err != nil {
		response.OkWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("成功", c)
	}

}

func Update(c *gin.Context) {

}

func Destroy(c *gin.Context) {

}
