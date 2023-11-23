package before

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis 连接配置
type RedisConnetionStu struct {
	RedisHost   string
	RedisPort   string
	RedisPass   string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

var redisClient *redis.Client

/**
 * 连接Mysql服务
 */
func ConnectionRedis(config *RedisConnetionStu) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:        config.RedisHost + ":" + config.RedisPort,
		Password:    config.RedisPass,
		DialTimeout: config.DialTimeout,
		ReadTimeout: config.ReadTimeout,
	})
}

// 连接测试
func RedisPing(ctx context.Context) (string, error) {
	return redisClient.Ping(ctx).Result()
}

// 关闭redis，回收资源时使用
func CloseRedis() {
	errMsg := recover()
	if errMsg != nil {
		fmt.Printf("[close.redis.error] : %v", errMsg)
	}

	if redisClient != nil {
		redisClient.Close()
	}
}
