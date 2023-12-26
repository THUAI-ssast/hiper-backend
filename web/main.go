package main

import (
	"hiper-backend/api"
	"hiper-backend/config"
	"hiper-backend/model"
	"hiper-backend/mq"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.InitConfig()
	model.InitDb()
	model.AutoMigrateDb()
	model.InitRedis()
	config.InitConfigAfter()

	if !viper.GetBool("is_debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	go mq.ListenMsgForMatchFinished(mq.Ctx_callback, "match_finished")

	api.ApiListenHttp()
	//api.ApiListenHttps()
}
