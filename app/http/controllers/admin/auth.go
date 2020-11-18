package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhimma/goin-web/app/service"
	"github.com/zhimma/goin-web/database/model"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
	"github.com/zhimma/goin-web/library/jwt"
	"net/http"
	"time"
)

type LoginData struct {
	Account  string `json:"account" form:"account" binding:"required" zh:"账号"`
	Password string `json:"password" form:"account" binding:"required" zh:"密码"`
}

type RegisterData struct {
	Account    string `json:"account" form:"account" binding:"required" zh:"账号"`
	Password   string `json:"password" form:"account" binding:"required" zh:"密码"`
	RePassword string `json:"RePassword" form:"RePassword" binding:"required,eqfield=Password" zh:"重复密码"`
}

// 登陆
func Login(c *gin.Context) {
	var loginParams LoginData
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
	adminInfo, err1 := service.AdminLogin(params)
	if err1 != nil {
		response.FailWithMessage(err1.Error(), c)
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
	tokenDetail, err2 := j.GenerateJwtToken(adminInfo.ID)
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err2,
		})
		return
	}
	err3 := service.CacheAdminUserToken(adminInfo.ID, tokenDetail)
	if err3 != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err3,
		})
		return
	}
	tokenData := map[string]string{
		"access_token":  tokenDetail.AccessToken,
		"refresh_token": tokenDetail.RefreshToken,
	}
	c.JSON(http.StatusOK, gin.H{
		"data": tokenData,
	})
	return
}

// 注册用户
func Register(c *gin.Context) {
	// 数据检验
	var u RegisterData
	if err := c.ShouldBind(&u); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			response.FailWithMessage(err.Error(), c)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		response.FailWithMessage(errs.Translate(globalInstance.Translator), c)
		return
	}

	// 判断账号是否已经存在
	checkMap := map[string]interface{}{
		"Account": u.Account,
	}
	if status, err := service.CheckAdminField(checkMap); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		if status == true {
			response.FailWithMessage("该账号已存在，请更换", c)
			return
		}
	}

	// 注册用户
	salt := helper.RandStringBytes(4)
	password, _ := helper.GenerateHashString(u.Password, salt)
	data := model.Admin{
		Account:     u.Account,
		Salt:        salt,
		Password:    password,
		Avatar:      "",
		Name:        u.Account,
		Phone:       u.Account,
		Status:      0,
		LastLoginIp: "",
		LastLoginAt: time.Now(),
		LoginTimes:  0,
	}
	if err := service.AdminRegister(&data); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("注册成功", c)

}
func Logout(c *gin.Context) {

}
