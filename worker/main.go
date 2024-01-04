package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/THUAI-ssast/hiper-backend/web/model"

	"github.com/THUAI-ssast/hiper-backend/worker/mq"
	"github.com/THUAI-ssast/hiper-backend/worker/repository"
	"github.com/THUAI-ssast/hiper-backend/worker/task"
)

func main() {
	InitConfig()
	model.InitDb()
	model.InitRedis()
	mq.InitMq()

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
			err := task.Build(repository.DomainType(taskType), id)
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

// InitConfig initializes the configuration of the application
func InitConfig() {
	viper.AutomaticEnv()
	// We can use `redis.host` instead of `REDIS_HOST`
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read remaining configs from file
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}

func getIDFromValues(values map[string]interface{}) uint {
	idInt, err := strconv.Atoi(values["id"].(string))
	if err != nil {
		log.Fatal(err)
	}
	return uint(idInt)
}
