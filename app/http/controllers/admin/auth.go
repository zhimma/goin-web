package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
	"github.com/zhimma/goin-web/library/jwt"
	"github.com/zhimma/goin-web/service"
	"net/http"
)

type RegisterData struct {
	Account  string `json:"account" form:"account" binding:"required" zh:"账号"`
	Password string `json:"password" form:"account" binding:"required" zh:"密码"`
}

func Login(c *gin.Context) {
	var loginParams RegisterData
	if err := c.ShouldBindJSON(&loginParams); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			response.FailWithMessage(err.Error(), c)
			return
		}
		errorMessageBag := helper.RemoveTopStruct(errs.Translate(globalInstance.Translator))
		response.FailWithMessage(errorMessageBag, c)
		return
	}

	params := &model.Admin{
		Account:  loginParams.Account,
		Password: loginParams.Password,
	}
	adminInfo, err := service.AdminLogin(params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if adminInfo.Status == 0 {
		response.FailWithMessage("账号已禁用或不存在", c)
		return
	}
	if helper.CompareHashString(adminInfo.Password, []byte(params.Password+adminInfo.Salt)) == false {
		response.FailWithMessage("密码错误，请检查后重新登陆", c)
		return
	}
	// 获取jwt结构体实例
	j := jwtLibrary.JWT{
		SigningKey:          []byte(globalInstance.BaseConfig.Jwt.JwtSecret),
		AccessTokenExpires:  globalInstance.BaseConfig.Jwt.JwtTtl,
		RefreshTokenExpires: globalInstance.BaseConfig.Jwt.JwtRefreshTtl,
	}
	tokenDetail, err := j.GenerateJwtToken(adminInfo.ID)
	service.CacheAdminUserToken(1, tokenDetail)
	tokenData := map[string]string{
		"access_token":  tokenDetail.AccessToken,
		"refresh_token": tokenDetail.RefreshToken,
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  tokenDetail,
		"data2": tokenData,
	})
	return
}

func Register(c *gin.Context) {
	var u RegisterData
	if err := c.ShouldBind(&u); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			"msg": errs.Translate(globalInstance.Translator),
		})
		return
	}
	// 保存入库等业务逻辑代码...

	c.JSON(http.StatusOK, "success")

	/*var params RegisterData
	if err := c.ShouldBindJSON(&params); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			"msg": errs.Translate(trans),
		})
		return
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		fmt.Println(params.Account)
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"name": "马雄飞"})*/
}
func Logout(c *gin.Context) {

}
