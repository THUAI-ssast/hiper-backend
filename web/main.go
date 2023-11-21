package main

import (
	"hiper-backend/api"
	"hiper-backend/config"
	"hiper-backend/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// TODO: Call init functions
	config.InitConfig()
	model.InitDB()
	model.InitRedis()

	if !viper.GetBool("is_debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	api.ApiListenHttp()
}
