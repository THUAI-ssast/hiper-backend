package mq

import (
	"context"
	"fmt"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/redis/go-redis/v9"
)

func SendBuildAIMsg(ctx context.Context, aiID uint) error {
	rdb := model.Rdb

	_, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "build",
		Values: map[string]interface{}{
			"type": "ai",
			"id":   fmt.Sprintf("%d", aiID),
		},
	}).Result()
	if err != nil {
		return err
	}

	return nil
}

func SendBuildGameLogicMsg(ctx context.Context, gameID uint) error {
	rdb := model.Rdb

	_, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "build",
		Values: map[string]interface{}{
			"type": "game_logic",
			"id":   fmt.Sprintf("%d", gameID),
		},
	}).Result()
	if err != nil {
		return err
	}

	return nil
}

func SendRunAutoMatchMsg(ctx context.Context, matchID uint) error {
	rdb := model.Rdb

	_, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "auto_match",
		Values: map[string]interface{}{
			"id": fmt.Sprintf("%d", matchID),
		},
	}).Result()
	if err != nil {
		return err
	}

	return nil
}
