package mq

import (
	"context"
	"fmt"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/redis/go-redis/v9"
)

func SendBuildAIMsg(ctx context.Context, aiID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Build",
		Body:  []byte(fmt.Sprintf("%d", aiID)),
		Type:  "AI",
	})
}

func SendBuildGameLogicMsg(ctx context.Context, gameID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Build",
		Body:  []byte(fmt.Sprintf("%d", gameID)),
		Type:  "Game",
	})
}

func SendBuildSdkMsg(ctx context.Context, sdkID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Build",
		Body:  []byte(fmt.Sprintf("%d", sdkID)),
		Type:  "SDK",
	})
}

func SendRunMatchMsg(ctx context.Context, matchID uint) error {
	args := &redis.XAddArgs{
		Stream: "manual_match",
		Values: map[string]interface{}{
			"id": fmt.Sprintf("%d", matchID),
		},
	}
	_, err := model.Rdb.XAdd(ctx, args).Result()
	return err
}
