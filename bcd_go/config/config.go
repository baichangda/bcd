package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

var redisAddress = "cdbai.cn:6379"
var redisPwd = "bcd5221527"

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() {

	if RedisClient == nil {
		env1 := os.Getenv("REDIS_HOST")
		env2 := os.Getenv("REDIS_PWD")
		if env1 != "" {
			redisAddress = env1
		}
		if env2 != "" {
			redisPwd = env2
		}
		RedisClient = redis.NewClient(&redis.Options{
			Addr:         redisAddress,
			Password:     redisPwd,
			DialTimeout:  10 * time.Second, // 设置连接超时
			ReadTimeout:  10 * time.Second, // 设置读取超时
			WriteTimeout: 10 * time.Second, // 设置写入超时
		})
	}
}
