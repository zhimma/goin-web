package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

const (
	ERROR   = -1
	SUCCESS = 0
)

func Result(statusCode int, code int, data interface{}, message string, c *gin.Context) {
	c.JSON(statusCode, Response{
		code,
		data,
		message,
	})
}

func Ok(c *gin.Context) {
	Result(http.StatusOK, SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(http.StatusOK, SUCCESS, data, "操作成功", c)
}
func Fail(c *gin.Context) {
	Result(http.StatusOK, ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, ERROR, map[string]interface{}{}, message, c)
}

func Unauthorized(message string, c *gin.Context) {
	Result(http.StatusUnauthorized, ERROR, map[string]interface{}{}, message, c)
}
