package config

import (
	"github.com/redis/go-redis/v9"
	"go_code/RPC/utils"
	"time"
)

func RedisConnect() {
	client := redis.NewClient(&redis.Options{
		Addr:        "redis-14520.c299.asia-northeast1-1.gce.cloud.redislabs.com:14520",
		Password:    "rPYdtUeiD5CeJSqcGZdoyHDd6Ou2uApa",
		DB:          0,
		DialTimeout: time.Second * 5,
	})
	utils.Redis = client
}
