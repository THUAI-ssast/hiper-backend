package repository

// a bunch of boring functions to update the state of a task

import (
	"github.com/THUAI-ssast/hiper-backend/web/model"
)

func StartBuildTask(taskType string, id uint) {
	switch taskType {
	case "game_logic":
		model.UpdateGameByID(id, map[string]interface{}{"game_logic_state": model.TaskStateRunning})
	case "ai":
		model.UpdateAiByID(id, map[string]interface{}{"task_state": model.TaskStateRunning})
	}
}

func EndBuildTask(taskType string, id uint, state model.TaskState, msg string) {
	switch taskType {
	case "game_logic":
		model.UpdateGameByID(id, map[string]interface{}{"game_logic_state": state, "game_logic_msg": msg})
	case "ai":
		model.UpdateAiByID(id, map[string]interface{}{"task_state": state, "task_msg": msg})
	}
}

func StartMatchTask(matchID uint) {
	model.UpdateMatchByID(matchID, map[string]interface{}{"state": model.TaskStateRunning})
}

func EndMatchTask(matchID uint, state model.TaskState) {
	model.UpdateMatchByID(matchID, map[string]interface{}{"state": state})
}

func StartGameLogicBuildDockerTask(gameID uint) {
	model.UpdateGameByID(gameID, map[string]interface{}{"game_logic_build_state": model.TaskStateRunning})
}

func EndGameLogicBuildDockerTask(gameID uint, state model.TaskState, msg string) {
	model.UpdateGameByID(gameID, map[string]interface{}{"game_logic_build_state": state, "game_logic_build_msg": msg})
}

// TODO: add more functions
