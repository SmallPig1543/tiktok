package cache

import (
	"github.com/redis/go-redis/v9"
	"tiktok/config"
)

var RedisClient *redis.Client

func LinkRedis() {
	conf := config.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost + ":" + conf.RedisPort,
		Password: conf.RedisPassword,
		DB:       conf.RedisDbName,
	})
	RedisClient = rdb
}
