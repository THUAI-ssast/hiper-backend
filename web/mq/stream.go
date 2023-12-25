package mq

import (
	"context"
	"hiper-backend/model"

	"github.com/redis/go-redis/v9"
)

func sendMsg(ctx context.Context, msg *Msg) error {
	return model.Rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: msg.Topic,
		MaxLen: 0,
		Approx: true,
		ID:     "*",
		Values: []interface{}{"body", msg.Body, "type", msg.Type},
	}).Err()
}

func SendByteMsg(ctx context.Context, topic string, body []byte, Type string) error {
	return model.Rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		MaxLen: 0,
		Approx: false,
		ID:     "*",
		Values: []interface{}{"body", body, "type", Type},
	}).Err()
}
