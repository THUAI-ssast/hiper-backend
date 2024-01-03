package repository

import (
	"github.com/THUAI-ssast/hiper-backend/web/model"
)

// UpdateBuildState updates the state of a build task
func UpdateBuildState(values map[string]interface{}, state model.TaskState) {
	// TODO: update build state
}

// UpdateMatchState updates the state of a match task
func UpdateMatchState(match_id uint, state model.TaskState) {
	// TODO: update match state
}
