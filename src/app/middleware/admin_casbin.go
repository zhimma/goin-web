package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/component"
	"github.com/zhimma/goin-web/database/model"
	"github.com/zhimma/goin-web/global/response"
	"net/http"
	"strings"
)

func AdminCasbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, exist := c.Get("adminInfo")
		if !exist {
			response.Abort(http.StatusInternalServerError, "用户登陆状态获取失败或获取UID出错", c)
			return
		}
		userInfo, ok := userData.(model.Admin)
		if ok {
			userInfo = userData.(model.Admin)
		} else {
			response.Abort(http.StatusInternalServerError, "用户登陆状态获取失败或获取UID出错", c)
			return
		}
		if userInfo.ID != 1 {
			url := c.Request.URL.RequestURI()
			method := c.Request.Method
			cas := component.Casbin()

			success, err := cas.Enforce("1", url, strings.ToUpper(method))

			if err != nil {
				response.Abort(http.StatusInternalServerError, "查询用户权限出错", c)
				return
			}
			if !success {
				response.Abort(http.StatusServiceUnavailable, "没有权限访问", c)
				return
			}
		}
		c.Next()

	}
}
