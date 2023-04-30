package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const redisAddress = "43.138.63.155:6379"
const redisPwd = "bcd"

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() {
	if RedisClient == nil {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:         redisAddress,
			Password:     redisPwd,
			DialTimeout:  10 * time.Second, // 设置连接超时
			ReadTimeout:  10 * time.Second, // 设置读取超时
			WriteTimeout: 10 * time.Second, // 设置写入超时
		})
	}
}
