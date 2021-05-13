package component

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/helper"
	"time"
)

type CacheStore interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expires time.Duration) error
	Delete(key string) error
	Exists(key string) error
	Replace(key string, value interface{}, expire time.Duration) error
}
type RedisStore struct {
	Client *redis.Client
}

func NewRedisCache() *RedisStore {
	return &RedisStore{Client: globalInstance.RedisClient}
}

func (c *RedisStore) Get(key string, value interface{}) error {
	err := c.Exists(key)
	if err != nil {
		return err
	}
	bytesData, _ := c.Client.Get(key).Bytes()
	return helper.Unserialize(bytesData, value) //反序列化操作
}

//保存数据
func (c *RedisStore) Set(key string, value interface{}, expires time.Duration) error {
	b, errs := helper.Serialize(value) //序列化操作，序列化可以保存对象
	if errs != nil {
		return errs
	}
	return c.Client.Set(key, b, expires).Err()
}

func (c *RedisStore) Delete(key string) error {
	err := c.Exists(key)
	if err != nil {
		return err
	}
	return c.Client.Del(key).Err()
}

// 判断key存在不存在
func (c *RedisStore) Exists(key string) error {
	err := c.Client.Exists(key).Err()
	if err == redis.Nil {
		return errors.New(fmt.Sprintf("the key %s not exists", key))
	}
	return nil
}
