package router

import (
	"errors"
	"github.com/zhimma/goin-web/internal/pkg/core"
	"github.com/zhimma/goin-web/internal/repository/mysql"
	"github.com/zhimma/goin-web/internal/repository/redis"
	"go.uber.org/zap"
)

type resource struct {
	logger *zap.Logger
	db     mysql.Repo
	cache  redis.Repo
	mux    core.Mux
}
type Server struct {
	Db    mysql.Repo
	Cache redis.Repo
	Mux   core.Mux
}

func NewHTTPServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	r := new(resource)
	r.logger = logger

	// 初始化数据库
	dbRepo, err := mysql.New()
	if err != nil {
		logger.Fatal("init database err", zap.Error(err))
	}
	r.db = dbRepo

	// 初始化cache
	cacheRepo, err := redis.New()
	if err != nil {
		logger.Fatal("init redis cache err", zap.Error(err))
	}
	r.cache = cacheRepo
	mux, err := core.New(logger, core.WithEnableCors(), core.WithEnableRate())
	if err != nil {
		panic(err)
	}
	r.mux = mux

	// 设置 API 路由
	setApiRouter(r)

	// 设置 Socket 路由
	// setSocketRouter(r)

	s := new(Server)
	s.Db = r.db
	s.Cache = r.cache
	s.Mux = mux

	return s, nil
}
