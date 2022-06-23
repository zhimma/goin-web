package main

import (
	"context"
	"github.com/zhimma/goin-web/config"
	"github.com/zhimma/goin-web/internal/router"
	"github.com/zhimma/goin-web/pkg/logger"
	"github.com/zhimma/goin-web/pkg/shutdown"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {

	// 初始化日志
	accessLogger, err := logger.NewJSONLogger(
		// 关闭console输出配置
		logger.WithDisableConsole(),
		// 添加额外字段配置
		logger.WithField("domain", "hello"),
		// 时间格式化
		logger.WithTimeLayout(config.CSTLayout),
		// 轮转切割日志
		logger.WithFileRotationP(config.ProjectAccessLogFile),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = accessLogger.Sync()
	}()
	// 初始化http服务
	s, err := router.NewHTTPServer(accessLogger)

	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    config.HttpServerPort,
		Handler: s.Mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			accessLogger.Fatal("http server startup err", zap.Error(err))
		}
	}()
	// 优雅关闭服务
	shutdown.NewHook().Close(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				accessLogger.Error("server shutdown err", zap.Error(err))
			}
		},

		func() {
			if s.Db != nil {
				if err := s.Db.DbRClose(); err != nil {
					accessLogger.Error("dbr close err", zap.Error(err))
				}
				if err := s.Db.DbWClose(); err != nil {
					accessLogger.Error("dbw close err", zap.Error(err))
				}
			}
		},

		func() {
			if s.Cache != nil {
				if err := s.Cache.Close(); err != nil {
					accessLogger.Error("cache close err", zap.Error(err))
				}
			}
		},
	)

}
