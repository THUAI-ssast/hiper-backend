package model

import (
	"errors"

	"gorm.io/gorm"
)

type Sdk struct {
	gorm.Model
	GameId uint `gorm:"index"`
	// ContestId is 0 if the SDK is in a game instead of any contest.
	ContestId uint `gorm:"index"`
	// Number is a unique identifier for each SDK within a game or contest.
	Number uint `gorm:"index"`

	Name   string `gorm:"unique"`
	Readme string

	BuildAi DockerTask `gorm:"embedded;embeddedPrefix:build_ai_"`
	RunAi   DockerTask `gorm:"embedded;embeddedPrefix:run_ai_"`
}

func (s *Sdk) BeforeCreate(tx *gorm.DB) (err error) {
	// Fill GameId from ContestId
	if s.ContestId != 0 && s.GameId == 0 {
		var gameId uint
		if err = tx.Model(&Contest{}).Select("game_id").First(&Contest{}, s.ContestId).Error; err != nil {
			return err
		}
		s.GameId = gameId
	}
	// Fill Number
	var maxNumber uint
	if err = tx.Model(&Sdk{}).Where("game_id = ? AND contest_id = ?", s.GameId, s.ContestId).Pluck("MAX(number)", &maxNumber).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		maxNumber = 0 // If no rows are found, set maxNumber to 0
	}
	s.Number = maxNumber + 1
	return nil
}

// TODO: add CRUD functions for sdk

// CRUD: Read

func GetSdks(filter map[string]interface{}, fields ...string) ([]Sdk, error) {
	// TODO: implement
	return []Sdk{}, nil
}
