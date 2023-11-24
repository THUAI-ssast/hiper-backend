package model

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var rdb *redis.Client

// InitRedis initializes the redis connection
func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
