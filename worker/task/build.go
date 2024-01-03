package task

import (
	"crypto/md5"
	"fmt"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/THUAI-ssast/hiper-backend/worker/repository"
)

func Build(domain repository.DomainType, id uint) (err error) {
	repository.StartBuildTask(domain, id)
	var taskState model.TaskState
	var msg string
	switch domain {
	case repository.GameLogicDomain:
		taskState, msg, err = buildGameLogic(id)
	case repository.AiDomain:
		taskState, msg, err = buildAI(id)
	}
	repository.EndBuildTask(domain, id, taskState, msg)
	return err
}

// 获取任务所需信息
// 起容器，执行任务
// 等待任务完成，获取任务输出，保存与修改相关信息

func buildGameLogic(gameID uint) {
	// TODO
}

func buildAI(aiID uint) (taskState model.TaskState, msg string, err error) {
	// TODO
}
