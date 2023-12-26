package model

import (
	"context"
	"fmt"
	"hiper-backend/mq"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var Rdb *redis.Client
var Ctx = context.Background()

// InitRedis initializes the redis connection
func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	go mq.ListenMsgForMatchFinished(mq.Ctx_callback, "match_finished")
}
