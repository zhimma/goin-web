package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/component"
	"github.com/zhimma/goin-web/global/response"
	"github.com/zhimma/goin-web/helper"
	"net/http"
	"strconv"
	"strings"
)

func AdminCasbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		managerRecord, errs := helper.GetCurrentManagerInfo(c)
		if errs != nil {
			response.Abort(http.StatusInternalServerError, errs.Error(), c)
			return
		}
		if managerRecord.IsSuper == 0 {
			url := c.Request.URL.RequestURI()
			method := c.Request.Method
			cas := component.Casbin()
			roleId := strconv.FormatInt(managerRecord.RoleId, 10)
			success, err := cas.Enforce(roleId, url, strings.ToUpper(method))
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
