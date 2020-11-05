package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type RegisterData struct {
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"account" binding:"required"`
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"name": "马雄飞"})
}

type SignUpParam struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

func Register(c *gin.Context) {
	var u SignUpParam
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
			"msg": errs.Translate(trans),
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
