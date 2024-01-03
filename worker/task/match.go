package task

import (
	"github.com/THUAI-ssast/hiper-backend/web/model"

	"github.com/THUAI-ssast/hiper-backend/worker/repository"
)

// TODO
func Match(matchID uint) (err error) {
	repository.StartMatchTask(matchID)
	// 获取任务所需信息

	// 起容器，执行任务

	// 等待任务完成，获取任务输出，保存与修改相关信息（含 在 match_result 消息队列中发送消息，如果是 公开对局 的话）

	repository.EndMatchTask(matchID, model.TaskStateFinished)
	return
}
