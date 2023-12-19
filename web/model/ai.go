package model

import (
	"errors"

	"gorm.io/gorm"
)

type Ai struct {
	gorm.Model
	GameId uint `gorm:"index"`
	// ContestId is 0 if the AI is in a game instead of any contest.
	ContestId uint `gorm:"index"`
	// Number is a unique identifier for each AI within a game or contest.
	Number uint `gorm:"index"`

	UserId uint `gorm:"index"`
	User   User `gorm:"foreignKey:UserId"`
	SdkId  uint
	Sdk    Sdk `gorm:"foreignKey:SdkId"`

	Note   string
	Status TaskStatus `gorm:"embedded;embeddedPrefix:task_"`
}

func (a *Ai) BeforeCreate(tx *gorm.DB) (err error) {
	// Fill GameId from ContestId
	if a.ContestId != 0 && a.GameId == 0 {
		var gameId uint
		if err = tx.Model(&Contest{}).Select("game_id").First(&Contest{}, a.ContestId).Error; err != nil {
			return err
		}
		a.GameId = gameId
	}
	// Fill Number
	var maxNumber uint
	if err = tx.Model(&Ai{}).Where("game_id = ? AND contest_id = ?", a.GameId, a.ContestId).Pluck("MAX(number)", &maxNumber).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		maxNumber = 0 // If no rows are found, set maxNumber to 0
	}
	a.Number = maxNumber + 1
	return nil
}

// TODO: add CRUD functions for ai

// CRUD: Read

func GetAis(query QueryParams) ([]Ai, int64, error) {
	// TODO: implement
	return []Ai{}, 0, nil
}
