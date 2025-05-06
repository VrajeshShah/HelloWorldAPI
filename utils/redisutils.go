package redisutils

import (
	"github.com/go-redis/redis"
)

type RedisObject struct {
	redisClient *redis.Client
	clientError error
}

func InitRedis() RedisObject {
	rValue := RedisObject{
		redisClient: redis.NewClient(&redis.Options{
			Addr: "192.168.1.102:6379",
		}),
	}
	_, err := rValue.redisClient.Ping().Result()
	rValue.clientError = err
	return rValue
}

func (redisObject *RedisObject) Get(key string) string {
	rValue := ""
	if redisObject.clientError != nil {
		return rValue
	}
	rValue, err := redisObject.redisClient.Get(key).Result()
	if err == redis.Nil {
		return rValue
	}
	return rValue
}

func (redisObject *RedisObject) Set(key string, value string) error {
	if redisObject.clientError != nil {
		return redisObject.clientError
	}
	rValue := redisObject.redisClient.Set(key, value, 0).Err()
	if rValue == nil {
		return nil
	}
	return rValue
}
