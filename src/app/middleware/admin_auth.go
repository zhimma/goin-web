package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/service/managerService"
	"github.com/zhimma/goin-web/global/response"
	jwtLibrary "github.com/zhimma/goin-web/library/jwt"
	"net/http"
	"time"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		message := "用户登陆失败:" + http.StatusText(http.StatusUnauthorized) + " -「%s」"
		// 从请求中获取token
		token := jwtLibrary.ExtractToken(c.Request)
		if token == "" {
			message = fmt.Sprintf(message, "MissingToken")
			response.Abort(http.StatusUnauthorized, message, c)
			return
		}

		jwt := jwtLibrary.NewJWT()
		tokenInfo, err := jwt.ParseJwtToken(token)
		if err != nil {
			message = fmt.Sprintf(message, "ParseTokenError:"+err.Error())
			response.Abort(http.StatusUnauthorized, message, c)
			return
		}

		// 查询redis中是否存在该token
		message = fmt.Sprintf(message, "Unauthorized")
		content, errs := managerService.AdminUserTokenCheck(tokenInfo)
		if errs != nil {
			response.Abort(http.StatusUnauthorized, message, c)
			return
		}
		if content == "" {
			response.Abort(http.StatusUnauthorized, message, c)
			return
		}
		passportResult, _ := managerService.GetManagerInfoFromCache(tokenInfo.IDENTIFIER)
		managerInfo := passportResult.ManagerInfo
		c.Set("managerInfo", managerInfo)
		c.Set("managerId", managerInfo.ID)
		//请求前获取当前时间
		nowTime := time.Now()
		//请求处理
		c.Next()

		//处理后获取消耗时间
		costTime := time.Since(nowTime)
		url := c.Request.URL.String()
		fmt.Printf("the request	 URL %s cost %v\n", url, costTime)
	}

}
