package main

import (
	"hiper-backend/api"
	"hiper-backend/config"
	"hiper-backend/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.InitConfig()
	model.InitDb()
	model.AutoMigrateDb()
	model.InitRedis()

	if !viper.GetBool("is_debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	//api.ApiListenHttp()
	api.ApiListenHttps()
}
