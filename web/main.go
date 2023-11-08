package main

import (
	"hiper-backend/api"
	"hiper-backend/config"
)

func main() {
	// TODO: Call init functions
	config.InitConfig()

	api.ApiListenHttp()

}
