package cache

import (
	"context"
	"forge/infra/configs"
	"forge/pkg/log/zlog"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	ctx         = context.Background() // Redis操作共享的context
)

func MustInitCache(config configs.IConfig) {
	err := initRedis(config)
	if err != nil {
		zlog.Errorf("初始化链接redis失败:%v", err.Error())
		panic(err)
	}
}
