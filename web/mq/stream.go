package mq

import (
	"context"
	"fmt"
	"hiper-backend/model"
	"strconv"
	"strings"

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

func ListenMsgForMatchFinished(ctx context.Context, topic string) (err error) {
	rdb := model.Rdb

	// 创建一个新的 XReadArgs
	args := &redis.XReadArgs{
		Streams: []string{topic, "$"},
		Block:   0,
	}

	for {
		// 读取消息
		streams, err := rdb.XRead(ctx, args).Result()
		if err != nil {
			return err
		}

		// 解析消息
		for _, stream := range streams {
			for _, message := range stream.Messages {
				body, ok := message.Values["body"].(string)
				if !ok {
					return fmt.Errorf("failed to parse body from message: %v", message)
				}

				_, ok = message.Values["type"].(string)
				if !ok {
					return fmt.Errorf("failed to parse type from message: %v", message)
				}

				// 解析 body
				parts := strings.Split(body, " ")
				if len(parts) < 2 {
					return fmt.Errorf("failed to parse matchID and replay from body: %v", body)
				}

				matchID, err := strconv.ParseUint(parts[0], 10, 32)
				if err != nil {
					return fmt.Errorf("failed to parse matchID from body: %v", body)
				}

				replay := parts[1]
				CallOnMatchFinished(uint(matchID), replay)
			}
		}
	}
}

func InitMq() {
	go ListenMsgForMatchFinished(Ctx_callback, "match_finished")
}

func InitGameMq(baseContestID uint) {
	SetCreateMatch(baseContestID)
	SetGetContestantsByRanking(baseContestID)
	SetUpdateContestant(baseContestID)
	SendBuildGameLogicMsg(Ctx_callback, baseContestID)
}
