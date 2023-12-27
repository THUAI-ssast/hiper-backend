package main

import (
	"log"
	"strconv"

	"github.com/THUAI-ssast/hiper-backend/web/model"

	"github.com/THUAI-ssast/hiper-backend/worker/mq"
	"github.com/THUAI-ssast/hiper-backend/worker/repository"
	"github.com/THUAI-ssast/hiper-backend/worker/task"
)

func main() {
	model.InitDb()
	model.InitRedis()

	for {
		stream, err := mq.GetTask()
		if err != nil {
			log.Println(err)
			continue
		}
		values := stream.Messages[0].Values

		switch stream.Stream {
		case "build":
			repository.UpdateBuildState(values, model.TaskStateRunning)
			task.Build(values)
			repository.UpdateBuildState(values, model.TaskStateFinished)
		case "manual_match", "auto_match":
			matchIDInt, err := strconv.Atoi(values["id"].(string))
			if err != nil {
				log.Fatal(err)
			}
			matchID := uint(matchIDInt)
			repository.UpdateMatchState(matchID, model.TaskStateRunning)
			task.Match(matchID)
			repository.UpdateMatchState(matchID, model.TaskStateFinished)
		}

		if err := mq.AckTask(stream); err != nil {
			log.Println(err)
		}
	}
}
