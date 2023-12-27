package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/THUAI-ssast/hiper-backend/api"
	"github.com/THUAI-ssast/hiper-backend/config"
	"github.com/THUAI-ssast/hiper-backend/model"
	"github.com/THUAI-ssast/hiper-backend/mq"
)

func main() {
	config.InitConfig()
	model.InitDb()
	model.AutoMigrateDb()
	model.InitRedis()
	mq.InitMq()
	config.InitConfigAfter()

	if !viper.GetBool("is_debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	api.ApiListenHttp()
	//api.ApiListenHttps()
}
