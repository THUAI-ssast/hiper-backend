package mq

import (
	"context"
	"fmt"
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
	return sendMsg(ctx, &Msg{
		Topic: "Run",
		Body:  []byte(fmt.Sprintf("%d", matchID)),
		Type:  "Match",
	})
}
