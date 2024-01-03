package task

import (
	"log"
	"strconv"
)

func Build(values map[string]interface{}) {
	idInt, err := strconv.Atoi(values["id"].(string))
	if err != nil {
		log.Fatal(err)
	}
	id := uint(idInt)

	switch values["type"] {
	case "game_logic":
		buildGameLogic(id)
	case "ai":
		buildAI(id)
	}
}

// 获取任务所需信息
// 起容器，执行任务
// 等待任务完成，获取任务输出，保存与修改相关信息

func buildGameLogic(gameID uint) {
	// TODO
}

func buildAI(aiID uint) {
	// TODO
}
