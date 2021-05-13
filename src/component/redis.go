package component

import (
	"fmt"
	"github.com/go-redis/redis"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
)

func RedisClient() {
	redisConfig := globalInstance.BaseConfig.Redis
	client := redis.NewClient(&redis.Options{
		// ork:            "",
		Addr: redisConfig.Addr,
		// Dialer:             nil,
		OnConnect: OnConnect,
		Password:  redisConfig.Password,
		DB:        redisConfig.Db,
		PoolSize:  100, // 连接池大小
		//MaxRetries:         0,
		//MinRetryBackoff:    0,
		//MaxRetryBackoff:    0,
		//DialTimeout:        0,
		//ReadTimeout:        0,
		//WriteTimeout:       0,
		//PoolSize:           0,
		//MinIdleConns:       0,
		//MaxConnAge:         0,
		//PoolTimeout:        0,
		//IdleTimeout:        0,
		//IdleCheckFrequency: 0,
		//TLSConfig:          nil,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		globalInstance.SystemLog.Error("redis connect ping failed, err:", zap.Any("err", err))
	} else {
		globalInstance.SystemLog.Info("redis connect ping response:", zap.Any("pong", pong))
		globalInstance.RedisClient = client
	}

}
func OnConnect(conn *redis.Conn) error {
	fmt.Println("redis on connect")
	globalInstance.RedisClientNum += 1
	return nil
}
