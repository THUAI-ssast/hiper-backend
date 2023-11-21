package model

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis() bool {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.ip"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("连接redis出错，错误信息：%v", err)
		return false
	}
	return true
}

func SetEX(key string, value interface{}, expiration int) bool {
	err := rdb.SetEX(ctx, key, value, time.Duration(expiration)*time.Hour).Err()
	if err != nil {
		fmt.Printf("redis set failed, err:%v\n", err)
		return false
	}
	return true
}
