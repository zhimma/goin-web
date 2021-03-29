package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Status     string      `json:"status"`
	HttpCode   int         `json:"httpCode"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
	Message    interface{} `json:"message"`
}

const (
	StatusSuccess     = "success"
	StatusError       = "error"
	StatusSuccessCode = 0
	StatusErrorCode   = -1
)

func Result(Status string, httpCode int, StatusCode int, data interface{}, message string, c *gin.Context) {
	c.JSON(httpCode, Response{
		Status,
		httpCode,
		StatusCode,
		data,
		message,
	})
}

func Abort(httpCode int, message string, c *gin.Context) {
	c.Abort()
	data := make([]interface{}, 0)
	Result(StatusError, httpCode, StatusErrorCode, data, message, c)
}

func Ok(c *gin.Context) {
	Result(StatusSuccess, http.StatusOK, StatusSuccessCode, []interface{}{}, "成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(StatusSuccess, http.StatusOK, StatusSuccessCode, []interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(StatusSuccess, http.StatusOK, StatusSuccessCode, data, "成功", c)
}
func Fail(c *gin.Context) {
	Result(StatusError, http.StatusOK, StatusErrorCode, []interface{}{}, "失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(StatusError, http.StatusOK, StatusErrorCode, []interface{}{}, message, c)
}

func ValidateFail(message string, c *gin.Context) {
	Result(StatusError, http.StatusBadRequest, StatusErrorCode, []interface{}{}, message, c)
}

func Unauthorized(message string, c *gin.Context) {
	Result(StatusError, http.StatusUnauthorized, http.StatusUnauthorized, []interface{}{}, message, c)
}
