package mq

import (
	"context"
	"errors"
	"os"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/redis/go-redis/v9"
)

var ctx context.Context
var rdb *redis.Client
var hostname string

func InitMq() {
	ctx = model.Ctx
	rdb = model.Rdb
	hostname = os.Getenv("HOSTNAME")
}

// GetTask returns a task from redis stream
func GetTask() (*redis.XStream, error) {
	t, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    "worker_group",
		Consumer: hostname,
		Streams:  []string{"build", "manual_match", "auto_match", ">"},
		Count:    1,
		Block:    0,
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(t) == 0 {
		return nil, errors.New("no task")
	}
	return &t[0], nil
}

// AckTask acknowledges a task
func AckTask(stream *redis.XStream) error {
	_, err := rdb.XAck(ctx, stream.Stream, "worker_group", stream.Messages[0].ID).Result()
	return err
}
