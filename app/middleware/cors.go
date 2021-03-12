/**

 * @Author: zhimma
 * @Description:
 * @File:  cors
 * @Date: 2020/11/17 10:09
 */
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/global/response"
	"net/http"
)

// 跨域中间件
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		origin := context.Request.Header.Get("Origin")
		mark := false
		cors := []string{
			"https://127.0.0.1",
			"http://127.0.0.1",
			"",
		}
		for _, v := range cors {
			if string(origin) == v {
				mark = true
				break
			}
		}
		if mark != true {
			response.Abort(http.StatusForbidden, -1, "Forbidden", context)
			return
		}

		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", origin)
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		context.Header("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X_Requested_With,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,Origin-Type,App-Identifier")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		if method == "OPTIONS" {
			response.Abort(http.StatusNoContent, 1, "Passed", context)
			return
		}
		// 处理请求
		context.Next()
	}
}
