package redis

import (
	"github.com/go-redis/redis/v8"
	"gobasic/ptc-Game/common/pkg/config"
	"golang.org/x/net/context"
	"time"
)

func New(config config.RedisConfig) (*redis.Client, error) {
	redisClient, err := newRedis(config)
	return redisClient, err
}

func newRedis(config config.RedisConfig) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  config.DialTimeOut,
		ReadTimeout:  config.ReadTimeOut,
		WriteTimeout: config.WriteTimeOut,
		PoolSize:     config.PoolSize,
		PoolTimeout:  config.PoolTimeOut,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
