package mq

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/redis/go-redis/v9"
)

func Startup() {
	rdb := model.Rdb

	streams := []string{"build", "manual_match", "auto_match"}
	group := "worker_group"
	for _, stream := range streams {
		err := rdb.XGroupCreateMkStream(Ctx_callback, stream, group, "0").Err()
		// ignore BUSYGROUP error
		if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
			panic(err)
		}
	}
}

func ListenMsgForMatchFinished() {
	rdb := model.Rdb

	for {
		streams, err := rdb.XRead(Ctx_callback, &redis.XReadArgs{
			Block:   0,
			Streams: []string{"match_result", "$"},
		}).Result()
		if err != nil {
			log.Println("Error reading stream: ", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				matchID, ok1 := message.Values["id"].(string)
				replay, ok2 := message.Values["replay"].(string)
				if ok1 && ok2 {
					matchIDUint, err := strconv.ParseUint(matchID, 10, 64)
					if err != nil {
						log.Println("Error parsing matchID: ", err)
						continue
					}
					CallOnMatchFinished(uint(matchIDUint), replay)
				}
			}
		}
	}
}
