package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/component"
	"github.com/zhimma/goin-web/global/response"
	"net/http"
	"strings"
)

func AdminCasbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		managerId, exist := c.Get("managerId")
		if !exist {
			response.Abort(http.StatusInternalServerError, "用户登陆状态获取失败或获取UID出错", c)
			return
		}
		if managerId != 1 {
			url := c.Request.URL.RequestURI()
			method := c.Request.Method
			cas := component.Casbin()
			role := 1
			success, err := cas.Enforce(role, url, strings.ToUpper(method))

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
