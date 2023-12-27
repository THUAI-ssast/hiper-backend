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
			match_id_str := values["match_id"].(string)
			match_id_int, err := strconv.Atoi(match_id_str)
			if err != nil {
				log.Fatal(err)
			}
			match_id := uint(match_id_int)
			repository.UpdateMatchState(match_id, model.TaskStateRunning)
			task.Match(match_id)
			repository.UpdateMatchState(match_id, model.TaskStateFinished)
		}

		if err := mq.AckTask(stream); err != nil {
			log.Println(err)
		}
	}
}
