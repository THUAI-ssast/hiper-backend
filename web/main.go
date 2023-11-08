package main

import (
	"flag"
	"hiper-backend/config"
	"hiper-backend/contest"
	"hiper-backend/game"
	"hiper-backend/user"
	"hiper-backend/web/api"
)

func main() {
	// TODO: Call init functions
	config.InitConfig()

	api.ApiListenHttp()

	var module string
	flag.StringVar(&module, "module", "", "assign run module")
	flag.Parse()
	switch module {
	case "game":
		game.Main()
	case "user":
		user.Main()
	case "contest":
		contest.Main()
	}
}
