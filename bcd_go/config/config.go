package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() {
	if RedisClient == nil {
		redisAddress := os.Getenv("REDIS_HOST")
		redisPwd := os.Getenv("REDIS_PWD")
		RedisClient = redis.NewClient(&redis.Options{
			Addr:         redisAddress,
			Password:     redisPwd,
			DialTimeout:  10 * time.Second, // 设置连接超时
			ReadTimeout:  10 * time.Second, // 设置读取超时
			WriteTimeout: 10 * time.Second, // 设置写入超时
		})
	}
}
