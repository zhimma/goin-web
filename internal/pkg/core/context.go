package core

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

const (
	_Alias            = "_alias_"
	_TraceName        = "_trace_"
	_LoggerName       = "_logger_"
	_BodyName         = "_body_"
	_PayloadName      = "_payload_"
	_GraphPayloadName = "_graph_payload_"
	_SessionUserInfo  = "_session_user_info"
	_AbortErrorName   = "_abort_error_"
	_IsRecordMetrics  = "_is_record_metrics_"
)

var contextPool = &sync.Pool{New: func() interface{} { return new(context) }}

func newContext(ctx *gin.Context) Context {
	context := contextPool.Get().(*context)
	context.ctx = ctx
	return context
}
func releaseContext(ctx Context) {
	c := ctx.(*context)
	c.ctx = nil
	contextPool.Put(c)
}

type Context interface {
	init()

	// ShouldBindQuery 反序列化 querystring
	// tag: `form:"xxx"` (注：不要写成 query)
	ShouldBindQuery(obj interface{}) error

	// ShouldBindPostForm 反序列化 postform (querystring会被忽略)
	// tag: `form:"xxx"`
	ShouldBindPostForm(obj interface{}) error

	// ShouldBindForm 同时反序列化 querystring 和 postform;
	// 当 querystring 和 postform 存在相同字段时，postform 优先使用。
	// tag: `form:"xxx"`
	ShouldBindFrom(obj interface{}) error

	// ShouldBindJSON 反序列化 postjson
	// tag: `json:"xxx"`
	ShouldBindJson(obj interface{}) error

	// ShouldBindURI 反序列化 path 参数(如路由路径为 /user/:name)
	// tag: `uri:"xxx"`
	ShouldBindURI(obj interface{}) error

	// Redirect 重定向
	Redirect(code int, location string)

	Logger() *zap.Logger
	setLogger(logger zap.Logger)

	Payload(payload interface{})
	getPayload() interface{}

	// AbortWithError 错误返回
	AbortWithError(err BusinessError)
	abortError() BusinessError

	Header() http.Header
	GetHeader(key string) string
	SetHeader(key, value string)

	Alias() string
	setAlias(path string)

	RequestInputParams() url.Values

	RequestPostFromParams() url.Values

	Request() *http.Request

	RawData() []byte

	Method() string

	Host() string

	Path() string

	URI() string

	ResponseWriter() gin.ResponseWriter
}

type context struct {
	ctx *gin.Context
}

func (c *context) init() {
	body, err := c.ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	c.ctx.Set(_BodyName, body)
	c.ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}

func (c *context) ShouldBindQuery(obj interface{}) error {
	return c.ctx.ShouldBindQuery(obj)
}

func (c *context) ShouldBindPostForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.FormPost)
}

func (c *context) ShouldBindFrom(obj interface{}) error {
	return c.ShouldBindFrom(obj)
}

func (c *context) ShouldBindJson(obj interface{}) error {
	return c.ctx.ShouldBindJSON(obj)
}

func (c *context) ShouldBindURI(obj interface{}) error {
	return c.ctx.ShouldBindUri(obj)
}

func (c *context) Redirect(code int, location string) {
	c.ctx.Redirect(code, location)
}

func (c *context) Logger() *zap.Logger {
	logger, ok := c.ctx.Get(_LoggerName)
	if !ok {
		return nil
	}
	return logger.(*zap.Logger)
}

func (c *context) setLogger(logger zap.Logger) {
	c.ctx.Set(_LoggerName, logger)
}

func (c *context) Payload(payload interface{}) {
	c.ctx.Set(_PayloadName, payload)
}

func (c *context) getPayload() interface{} {
	if payload, ok := c.ctx.Get(_PayloadName); ok != false {
		return payload
	}
	return nil
}

func (c *context) AbortWithError(err BusinessError) {
	if err != nil {
		httpCode := err.HTTPCode()
		if httpCode == 0 {
			httpCode = http.StatusInternalServerError
		}

		c.ctx.AbortWithStatus(httpCode)
		c.ctx.Set(_AbortErrorName, err)
	}
}

func (c *context) abortError() BusinessError {
	err, _ := c.ctx.Get(_AbortErrorName)
	return err.(BusinessError)
}

func (c *context) Header() http.Header {
	headers := c.ctx.Request.Header

	clone := make(http.Header, len(headers))
	// 把请求的header 重新给header
	for k, v := range headers {
		value := make([]string, len(v))
		copy(value, v)
		clone[k] = value
	}
	return clone
}

func (c *context) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}

func (c *context) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

func (c *context) Alias() string {
	panic("implement me")
}

func (c *context) setAlias(path string) {
	panic("implement me")
}

func (c *context) RequestInputParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.Form
}

func (c *context) RequestPostFromParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.PostForm
}

func (c *context) Request() *http.Request {
	return c.Request()
}

func (c *context) RawData() []byte {
	body, err := c.ctx.Get(_BodyName)
	if !err {
		return nil
	}
	return body.([]byte)
}

func (c *context) Method() string {
	return c.ctx.Request.Method
}

func (c *context) Host() string {
	return c.ctx.Request.Host
}

func (c *context) Path() string {
	return c.ctx.Request.URL.Path
}

func (c *context) URI() string {
	uri, _ := url.QueryUnescape(c.ctx.Request.URL.RequestURI())
	return uri
}

func (c *context) ResponseWriter() gin.ResponseWriter {
	return c.ctx.Writer
}

var _ Context = (*context)(nil)
