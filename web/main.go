package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/THUAI-ssast/hiper-backend/web/api"
	"github.com/THUAI-ssast/hiper-backend/web/config"
	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/THUAI-ssast/hiper-backend/web/mq"
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

	//TODO:DELETE!
	go mq.WarnCode()
	//TODO:DELETE!

	api.ApiListenHttp()
	//api.ApiListenHttps()
}
