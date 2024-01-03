package task

import (
	"crypto/md5"
	"fmt"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/THUAI-ssast/hiper-backend/worker/repository"
)

func Build(taskType string, id uint) (err error) {
	repository.StartBuildTask(taskType, id)
	var taskState model.TaskState
	var msg string
	switch taskType {
	case "game_logic":
		taskState, msg, err = buildGameLogic(id)
	case "ai":
		taskState, msg, err = buildAI(id)
	}
	repository.EndBuildTask(taskType, id, taskState, msg)
	return err
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
