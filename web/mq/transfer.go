package mq

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
)

func SendBuildAIMsg(ctx context.Context, aiID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Build",
		Body:  []byte(fmt.Sprintf("%d", aiID)),
		Type:  "AI",
	})
}

// Game和contest都会这样发送
func SendBuildGameMsg(ctx context.Context, gameID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Build",
		Body:  []byte(fmt.Sprintf("%d", gameID)),
		Type:  "Game",
	})
}

func SendBuildContestMsg(ctx context.Context, contestID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Build",
		Body:  []byte(fmt.Sprintf("%d", contestID)),
		Type:  "Contest",
	})
}

func SendBuildSdkMsg(ctx context.Context, sdkID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Build",
		Body:  []byte(fmt.Sprintf("%d", sdkID)),
		Type:  "SDK",
	})
}

func SendChangeGameMsg(ctx context.Context, gameID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Change",
		Body:  []byte(fmt.Sprintf("%d", gameID)),
		Type:  "Game",
	})
}

func SendBuildMatchMsg(ctx context.Context, matchID uint, extraInfo map[string]interface{}) error {
	// 将 extraInfo 转换为 JSON 格式的字符串
	extraInfoBytes, err := json.Marshal(extraInfo)
	if err != nil {
		return err
	}
	extraInfoStr := string(extraInfoBytes)

	// 将 matchID 和 extraInfoStr 添加到 Body 字段中
	body := fmt.Sprintf("%d,%s", matchID, extraInfoStr)

	return sendMsg(ctx, &Msg{
		Topic: "Build",
		Body:  []byte(body),
		Type:  "Match",
	})
}

func SendChangeMatchMsg(ctx context.Context, matchID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Change",
		Body:  []byte(fmt.Sprintf("%d", matchID)),
		Type:  "Match",
	})
}

func SendRunMatchMsg(ctx context.Context, matchID uint) error {
	return sendMsg(ctx, &Msg{
		Topic: "Run",
		Body:  []byte(fmt.Sprintf("%d", matchID)),
		Type:  "Match",
	})
}

func SendAIIDsMsg(ctx context.Context, aiIDs []uint) {
	buf := new(bytes.Buffer)

	for _, id := range aiIDs {
		err := binary.Write(buf, binary.LittleEndian, uint64(id))
		if err != nil {
			log.Fatalf("Failed to encode id: %v", err)
		}
	}

	data := buf.Bytes()

	SendByteMsg(ctx, "Run", data, "AI")
}
