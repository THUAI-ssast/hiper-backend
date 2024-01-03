package repository

// a bunch of boring functions to update the state of a task

import (
	"fmt"

	"github.com/THUAI-ssast/hiper-backend/web/model"
)

func StartBuildTask(domain DomainType, id uint) {
	switch domain {
	case GameLogicDomain:
		model.UpdateGameByID(id, map[string]interface{}{"game_logic_state": model.TaskStateRunning})
	case AiDomain:
		model.UpdateAiByID(id, map[string]interface{}{"task_state": model.TaskStateRunning})
	}
}

func EndBuildTask(domain DomainType, id uint, state model.TaskState, msg string) {
	switch domain {
	case GameLogicDomain:
		model.UpdateGameByID(id, map[string]interface{}{"game_logic_state": state, "game_logic_msg": msg})
	case AiDomain:
		model.UpdateAiByID(id, map[string]interface{}{"task_state": state, "task_msg": msg})
	}
}

func StartMatchTask(matchID uint) {
	model.UpdateMatchByID(matchID, map[string]interface{}{"state": model.TaskStateRunning})
}

func EndMatchTask(matchID uint, state model.TaskState) {
	model.UpdateMatchByID(matchID, map[string]interface{}{"state": state})
}

// domainType: game_logic, ai
// operationType: build, run
func StartBuildImageTask(domain DomainType, operation OperationType, id uint) {
	switch domain {
	case GameLogicDomain:
		field := fmt.Sprintf("game_logic_%s_state", operation)
		model.UpdateGameByID(id, map[string]interface{}{field: model.TaskStateRunning})
	case AiDomain:
		field := fmt.Sprintf("%s_ai_state", operation)
		model.UpdateSdkByID(id, map[string]interface{}{field: model.TaskStateRunning})
	}
}

func EndBuildImageTask(domain DomainType, operation OperationType, id uint, state model.TaskState, msg string) {
	// model.UpdateGameByID(gameID, map[string]interface{}{"game_logic_build_state": state, "game_logic_build_msg": msg})
	switch domain {
	case "game_logic":
		fieldState := fmt.Sprintf("game_logic_%s_state", operation)
		fieldMsg := fmt.Sprintf("game_logic_%s_msg", operation)
		model.UpdateGameByID(id, map[string]interface{}{fieldState: state, fieldMsg: msg})
	case "ai":
		fieldState := fmt.Sprintf("%s_ai_state", operation)
		fieldMsg := fmt.Sprintf("%s_ai_msg", operation)
		model.UpdateSdkByID(id, map[string]interface{}{fieldState: state, fieldMsg: msg})
	}
}

// TODO: add more functions
