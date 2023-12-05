package model

import (
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	BaseContest

	Admins []User `gorm:"many2many:game_admins;"`

	// TODO:
	// game assets
}

type GamePrivilege string

const (
	GamePrivilegeAdmin      GamePrivilege = "admin"
	GamePrivilegeRegistered GamePrivilege = "registered"
)

func CreateGame(game Game) error {
	return db.Create(&game).Error
}

// TODO
// GetGames returns all games
// readme is truncated to 100 characters
func GetGames() ([]Game, error) {
	return nil, nil
}

// TODO
// GetGameById returns game by id
func GetGameById(id uint) (Game, error) {
	return Game{}, nil
}

// TODO: add CRUD functions for game
