package cache

import (
	"context"
	"fmt"
	"forge/infra/configs"
	"forge/pkg/log/zlog"
	"github.com/go-redis/redis/v8"
)

const (
	redisAddr = "%s:%d"
)

func initRedis(config configs.IConfig) error {
	redisConfig := config.GetRedisConfig()
	if !config.GetRedisConfig().Enable {
		zlog.Warnf("不使用Redis模式")
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Network:            "",
		Addr:               fmt.Sprintf(redisAddr, redisConfig.Host, redisConfig.Port),
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           redisConfig.Password,
		DB:                 redisConfig.DB,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolFIFO:           false,
		PoolSize:           1000,
		MinIdleConns:       1,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		zlog.Errorf("redis无法链接 %v", err)
		return err
	}
	redisClient = client
	return nil
}
