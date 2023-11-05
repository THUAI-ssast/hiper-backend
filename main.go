package main

import (
	"flag"
	"fmt"
	"hiper-backend/contest"
	"hiper-backend/game"
	"hiper-backend/user"
)

func main() {
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
	fmt.Println("Done!")
}
