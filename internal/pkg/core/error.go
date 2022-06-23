package core

import "github.com/pkg/errors"

type BusinessError interface {
	i()

	// WithError 设置错误信息
	WithError(err error) BusinessError

	// withAlert 设置告警通知
	WithAlert() BusinessError

	// StatusCode 获取业务码
	StatusCode() int

	// httpCode 获取http状态吗
	HTTPCode() int

	// message 获取错误信息描述
	Message() string

	// StackError 获取带堆栈的错误信息
	StackError() error

	// isAlert 是否开启警告通知
	IsAlert() bool
}

type businessError struct {
	httpCode   int    // http状态吗
	statusCode int    // 业务吗
	message    string // 错误描述
	stackError error  // 含有堆栈信息的错误
	isAlert    bool   // 是否告警
}

func (e *businessError) i() {}

func Error(httpCode, statusCode int, message string) BusinessError {
	return &businessError{
		httpCode:   httpCode,
		statusCode: statusCode,
		message:    message,
		isAlert:    false,
	}
}
func (e *businessError) WithError(err error) BusinessError {
	e.stackError = errors.WithStack(err)
	return e
}

func (e *businessError) WithAlert() BusinessError {
	e.isAlert = true
	return e
}

func (e *businessError) StatusCode() int {
	return e.statusCode
}

func (e *businessError) HTTPCode() int {
	return e.httpCode
}

func (e *businessError) Message() string {
	return e.message
}

func (e *businessError) IsAlert() bool {
	return e.isAlert
}

func (e *businessError) StackError() error {
	return e.stackError
}
