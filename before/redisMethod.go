package before

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var appCtx = context.TODO()

/**
 * 查询所有key
 */
func RdbKeys() ([]string, error) {
	val, err := redisClient.Keys(appCtx, "*").Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return val, nil
}

/**
 * 设置string类型的数据
 * @Param key string 键名称
 * @Param value interface{} 值内容，传递的值都会作为字符串保存
 * @Param ttl time.Duration 过期时间，例如：time.Second * 1
 * @Return 为 nil 时表示添加成功
 */
func RdbSet(key string, val interface{}, ttl time.Duration) error {
	op := redisClient.Set(appCtx, key, val, ttl)
	return op.Err()
}

/**
 * 获取string类型数据
 * @Param key string 键名称
 * @Return 为 nil 时表示获取成功
 */
func RdbGet(key string) (string, error) {
	val, err := redisClient.Get(appCtx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return val, nil
}

/**
 * 添加Set类型的数据
 * @Param key string 键名称
 * @Param members 需要保存到集合中的值
 */
func RdbSAdd(key string, members ...interface{}) error {
	op := redisClient.SAdd(appCtx, key, members)
	return op.Err()
}

/**
 * 删除无序集合中的元素
 * @Param key string 键名称
 * @Param members 需要删除的值
 */
func RdbSRem(key string, members ...interface{}) (int64, error) {
	return redisClient.SRem(appCtx, key, members).Result()
}

/**
 * 返回无序集合中的个数
 * @Param key string 键名称
 */
func RdbSCard(key string) (int64, error) {
	num, err := redisClient.SCard(appCtx, key).Result()
	return num, err
}

/**
 * 判断某个值是否存在某个无序集合中
 * @Param key string 键名称
 * @Param val interface{} 需要判断的值
 */
func RdbSIsMember(key string, val interface{}) (bool, error) {
	return redisClient.SIsMember(appCtx, key, val).Result()
}

/**
 * 返回无序集合中的所有元素
 * @Param key string 键名称
 * @Return 返回的值都是以字符串保存的
 */
func RdbSMembers(key string) ([]string, error) {
	return redisClient.SMembers(appCtx, key).Result()
}

/**
 * 删除数据
 * @Param key string 需要删除的键的名称
 */
func RdbDel(key string) error {
	return redisClient.Del(appCtx, key).Err()
}

/**
 * 发布消息
 * @Param ch string 消息对应的通道的名称
 * @Param msg string 消息内容
 */
func RdbPublish(ctx context.Context, ch, msg string) error {
	return redisClient.Publish(ctx, ch, msg).Err()
}

/**
 * 订阅消息
 * @Param ch string 消息对应的通道的名称
 */
func RdbSubscribe(ctx context.Context, ch string) (string, error) {
	sub := redisClient.Subscribe(ctx, ch)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}
	return msg.Payload, nil
}
