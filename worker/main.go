package main

import (
	"log"
	"strconv"

	"github.com/THUAI-ssast/hiper-backend/web/model"

	"github.com/THUAI-ssast/hiper-backend/worker/mq"
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
			taskType := values["type"].(string)
			id := getIDFromValues(values)
			err := task.Build(taskType, id)
			if err != nil {
				log.Println(err)
			}
		case "manual_match", "auto_match":
			matchID := getIDFromValues(values)
			err := task.Match(matchID)
			if err != nil {
				log.Println(err)
			}
		}

		if err := mq.AckTask(stream); err != nil {
			log.Println(err)
		}
	}
}

func getIDFromValues(values map[string]interface{}) uint {
	idInt, err := strconv.Atoi(values["id"].(string))
	if err != nil {
		log.Fatal(err)
	}
	return uint(idInt)
}
