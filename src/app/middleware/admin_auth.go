package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/service/managerService"
	"github.com/zhimma/goin-web/database/model"
	"github.com/zhimma/goin-web/global/response"
	jwtLibrary "github.com/zhimma/goin-web/library/jwt"
	"net/http"
	"time"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		message := "用户登陆失败:" + http.StatusText(http.StatusUnauthorized) + " -「%s」"

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
		originUserData, _ := managerService.GetManagerInfoFromCache(tokenInfo.IDENTIFIER)
		managerInfo := originUserData["managerInfo"]
		userInfo, ok := managerInfo.(model.Manager)
		if ok {
			userInfo = managerInfo.(model.Manager)
		} else {
			response.Abort(http.StatusInternalServerError, "用户登陆状态获取失败或获取UID出错", c)
			return
		}
		c.Set("managerInfo", userInfo)
		c.Set("managerId", userInfo.ID)

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
