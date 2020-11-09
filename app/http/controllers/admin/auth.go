package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/library/jwt"
	"github.com/zhimma/goin-web/service"
	"net/http"
)

type RegisterData struct {
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"account" binding:"required"`
}

func Login(c *gin.Context) {
	// 获取jwt结构体实例
	j := jwtLibrary.JWT{
		SigningKey:          []byte(globalInstance.BaseConfig.Jwt.JwtSecret),
		AccessTokenExpires:  globalInstance.BaseConfig.Jwt.JwtTtl,
		RefreshTokenExpires: globalInstance.BaseConfig.Jwt.JwtRefreshTtl,
	}
	tokenDetail, err := j.GenerateJwtToken(1)

	service.CacheUserToken(1, tokenDetail)
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
