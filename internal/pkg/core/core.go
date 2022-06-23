package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/zhimma/goin-web/config"
	"github.com/zhimma/goin-web/internal/code"
	"github.com/zhimma/goin-web/pkg/env"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"net/http"
	"runtime/debug"
	"time"
)

type HandlerFunc func(c Context)

// RouterGroup 包装gin的RouterGroup
type RouterGroup interface {
	Group(string, ...HandlerFunc) RouterGroup
	IRoutes
}

func (r *router) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	group := r.group.Group(relativePath, wrapHandlers(handlers...)...)
	return &router{group: group}
}

// IRoutes 包装gin的IRoutes
type IRoutes interface {
	Any(string, ...HandlerFunc)
	GET(string, ...HandlerFunc)
	POST(string, ...HandlerFunc)
	DELETE(string, ...HandlerFunc)
	PATCH(string, ...HandlerFunc)
	PUT(string, ...HandlerFunc)
	OPTIONS(string, ...HandlerFunc)
	HEAD(string, ...HandlerFunc)
}

type router struct {
	group *gin.RouterGroup
}

func (r *router) Any(relativePath string, handlers ...HandlerFunc) {
	r.group.Any(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) GET(relativePath string, handlers ...HandlerFunc) {
	r.group.GET(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) POST(relativePath string, handlers ...HandlerFunc) {
	r.group.POST(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) DELETE(relativePath string, handlers ...HandlerFunc) {
	r.group.DELETE(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) PATCH(relativePath string, handlers ...HandlerFunc) {
	r.group.PATCH(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) PUT(relativePath string, handlers ...HandlerFunc) {
	r.group.PUT(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	r.group.OPTIONS(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) HEAD(relativePath string, handlers ...HandlerFunc) {
	r.group.HEAD(relativePath, wrapHandlers(handlers...)...)
}

var _ IRoutes = (*router)(nil)

func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		funcs[i] = func(context *gin.Context) {
			ctx := newContext(context)
			defer releaseContext(ctx)
			handler(ctx)
		}
	}
	return funcs
}

// mux builder
type Mux interface {
	http.Handler
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

type mux struct {
	engine *gin.Engine
}

var _ Mux = (*mux)(nil)

func (m *mux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.engine.ServeHTTP(writer, request)
}

func (m *mux) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	return &router{group: m.engine.Group(relativePath, wrapHandlers(handlers...)...)}
}

// options 配置管理
type option struct {
	enableCors bool
	enableRate bool
}
type Options func(*option)

func WithEnableCors() Options {
	return func(o *option) {
		o.enableCors = true
	}
}

func WithEnableRate() Options {
	return func(o *option) {
		o.enableRate = true
	}
}

func New(logger *zap.Logger, options ...Options) (Mux, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	gin.SetMode(gin.ReleaseMode)
	mux := &mux{
		engine: gin.New(),
	}

	opt := new(option)
	for _, f := range options {
		f(opt)
	}

	if opt.enableCors {
		mux.engine.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:     []string{"*"},
			AllowCredentials:   true,
			OptionsPassthrough: true,
		}))
	}

	mux.engine.Use(func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", string(debug.Stack())))
			}
		}()
		ctx.Next()
	})

	mux.engine.Use(func(ctx *gin.Context) {
		if ctx.Writer.Status() == http.StatusNotFound {
			return
		}
		// ts := time.Now()

		context := newContext(ctx)
		defer releaseContext(context)

		defer func() {
			var (
				response     interface{}
				statusCode   int
				errorMessage string
				abortErr     error
			)

			// 源发生错误，处理错误 返回
			if ctx.IsAborted() {
				for i := range ctx.Errors {
					multierr.AppendInto(&abortErr, ctx.Errors[i])
				}

				if err := context.abortError(); err != nil {
					// TODO 发送消息通知

					multierr.AppendInto(&abortErr, err.StackError())

					statusCode = err.StatusCode()
					errorMessage = err.Message()
					response = &code.Failure{
						StatusCode: statusCode,
						Message:    errorMessage,
					}
					ctx.JSON(err.HTTPCode(), response)
				}
			}

			// 源正确返回
			response = context.getPayload()
			if response != nil {
				ctx.JSON(http.StatusOK, response)
			}

			// TODO记录接口指标

			// TODO记录请求日志

		}()
		// 放行 到下一个
		ctx.Next()
	})

	if opt.enableRate {
		limiter := rate.NewLimiter(rate.Every(time.Second*1), config.MaxRequestsPerSecond)

		mux.engine.Use(func(ctx *gin.Context) {
			context := newContext(ctx)
			defer releaseContext(context)

			if !limiter.Allow() {
				context.AbortWithError(Error(
					http.StatusTooManyRequests,
					code.TooManyRequests,
					code.Text(code.TooManyRequests),
				))
				return
			}
			ctx.Next()
		})
	}

	system := mux.Group("/system")
	{
		system.GET("/health", func(ctx Context) {
			response := &struct {
				Timestamp   time.Time `json:"timestamp"`
				Environment string    `json:"environment"`
				Host        string    `json:"host"`
				Status      string    `json:"status"`
			}{
				Timestamp:   time.Now(),
				Environment: env.NowEnv().Value(),
				Host:        ctx.Host(),
				Status:      "ok",
			}
			ctx.Payload(response)
		})
	}
	return mux, nil
}
