package initiliza

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"order_bff/global"
)

// InitRedis 初始化redis
func InitRedis() {

	global.RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "",
	})

	if err := global.RedisClient.Ping(context.Background()).Err(); err != nil {
		zap.S().Info("连接redis失败")
		return
	}
}
