package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/service/adminAuthService"
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
		if content, err := adminAuthService.AdminUserTokenCheck(tokenInfo); err != nil {
			response.Abort(http.StatusUnauthorized, message, c)
			return
		} else {
			if content == "" {
				response.Abort(http.StatusUnauthorized, message, c)
				return
			}
		}
		c.Set("UID", tokenInfo.IDENTIFIER)

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