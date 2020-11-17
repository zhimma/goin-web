package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	jwtLibrary "github.com/zhimma/goin-web/library/jwt"
	"github.com/zhimma/goin-web/service"
	"net/http"
	"time"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		message := "用户登陆失败:" + http.StatusText(http.StatusUnauthorized) + " -「%s」"
		message = fmt.Sprintf(message, "Unauthorized")

		token := jwtLibrary.ExtractToken(c.Request)
		globalInstance.SystemLog.Info("用户登陆,token:" + token)
		if token == "" {
			message = fmt.Sprintf(message, "MissingToken")
			response.Abort(http.StatusUnauthorized, -1, message, message, c)
		}

		jwt := jwtLibrary.NewJWT()
		tokenInfo, err := jwt.ParseJwtToken(token)
		if err != nil {
			message = fmt.Sprintf(message, "ParseTokenError:"+err.Error())
			response.Abort(http.StatusUnauthorized, -1, message, message, c)
		}

		// 查询redis中是否存在该token
		if content, err := service.AdminUserTokenCheck(tokenInfo); err != nil {
			response.Abort(http.StatusUnauthorized, -1, message, message, c)
		} else {
			if content == "" {
				response.Abort(http.StatusUnauthorized, -1, message, message, c)
			}
		}
		c.Set("UUID", tokenInfo.UUID)
		c.Set("UID", tokenInfo.UID)

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
