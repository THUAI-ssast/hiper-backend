package repository

import (
	"strconv"

	"github.com/THUAI-ssast/hiper-backend/web/model"
)

// UpdateBuildState updates the state of a build task
func UpdateBuildState(values map[string]interface{}, state model.TaskState) {
	idInt, _ := strconv.Atoi(values["id"].(string))
	id := uint(idInt)

	switch values["type"] {
	case "game_logic":
		model.UpdateGameByID(id, map[string]interface{}{"game_logic_build_state": state})
	case "ai":
		model.UpdateAiByID(id, map[string]interface{}{"task_state": state})
	}
}

// UpdateMatchState updates the state of a match task
func UpdateMatchState(match_id uint, state model.TaskState) {
	model.UpdateMatchByID(match_id, map[string]interface{}{"state": state})
}
