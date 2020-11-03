package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求前获取当前时间
		nowTime := time.Now()
		fmt.Println("调用中间件方法")
		//请求处理
		c.Next()

		//处理后获取消耗时间
		costTime := time.Since(nowTime)
		url := c.Request.URL.String()
		fmt.Printf("the request URL %s cost %v\n", url, costTime)
	}
}
