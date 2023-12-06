package model

import (
	"errors"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model
	GameId uint `gorm:"index"`
	// ContestId is 0 if the match is in a game instead of any contest.
	ContestId uint `gorm:"index"`
	// Number is a unique identifier for each match within a game or contest.
	Number uint `gorm:"index"`

	Players []Ai `gorm:"many2many:match_ais;"`
	State   TaskState
	Tag     string

	Logs   []string `gorm:"serializer:json"`
	Replay string
	Scores []int `gorm:"serializer:json"`
}

func (m *Match) BeforeCreate(tx *gorm.DB) (err error) {
	// Fill GameId from ContestId
	if m.ContestId != 0 && m.GameId == 0 {
		var gameId uint
		if err = tx.Model(&Contest{}).Select("game_id").First(&Contest{}, m.ContestId).Error; err != nil {
			return err
		}
		m.GameId = gameId
	}
	// Fill Number
	var maxNumber uint
	if err = tx.Model(&Match{}).Where("game_id = ? AND contest_id = ?", m.GameId, m.ContestId).Pluck("MAX(number)", &maxNumber).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		maxNumber = 0 // If no rows are found, set maxNumber to 0
	}
	m.Number = maxNumber + 1
	return nil
}

// TODO: add CRUD functions for match
