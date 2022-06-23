package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/zhimma/goin-web/config"
	"time"
)

type option struct {
}

type Option func(option *option)

func newOption() *option {
	return &option{}
}

type Repo interface {
	i()
	Set(key, value string, ttl time.Duration, option ...Option) error
	Get(key string, option ...Option) (string, error)
	Del(key string, option ...Option) bool
	TTL(key string) (time.Duration, error)
	Expire(key string, ttl time.Duration) bool
	Exists(key string) bool
	Incr(key string, option ...Option) int64
	Close() error
	Version() string
}

type cacheRepo struct {
	client *redis.Client
}

func (c cacheRepo) i() {
	panic("implement me")
}

func (c *cacheRepo) Set(key, value string, ttl time.Duration, option ...Option) error {
	panic("implement me")
}

func (c *cacheRepo) Get(key string, option ...Option) (string, error) {
	panic("implement me")
}

func (c *cacheRepo) Del(key string, option ...Option) bool {
	panic("implement me")
}

func (c *cacheRepo) TTL(key string) (time.Duration, error) {
	panic("implement me")
}

func (c *cacheRepo) Expire(key string, ttl time.Duration) bool {
	panic("implement me")
}

func (c *cacheRepo) Exists(key string) bool {
	panic("implement me")
}

func (c *cacheRepo) Incr(key string, option ...Option) int64 {
	panic("implement me")
}

func (c *cacheRepo) Close() error {
	panic("implement me")
}

func (c *cacheRepo) Version() string {
	panic("implement me")
}

var _ Repo = (*cacheRepo)(nil)

func New() (Repo, error) {
	client, err := redisConnect()
	if err != nil {
		return nil, err
	}
	return &cacheRepo{client: client}, nil
}

func redisConnect() (*redis.Client, error) {
	cfg := config.Get().Redis
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Pass,
		DB:           cfg.Db,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrap(err, "ping redis err")
	}
	return client, nil
}
