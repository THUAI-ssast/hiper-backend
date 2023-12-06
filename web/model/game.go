package model

import (
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	BaseContest

	Admins []User `gorm:"many2many:game_admins;"`

	// game assets
	GameLogic   GameLogic   `gorm:"embedded;embeddedPrefix:game_logic_"`
	MatchDetail MatchDetail `gorm:"embedded;embeddedPrefix:match_detail_"`
}

type GamePrivilege string

const (
	GamePrivilegeAdmin      GamePrivilege = "admin"
	GamePrivilegeRegistered GamePrivilege = "registered"
)

type GameLogic struct {
	Build  DockerTask `gorm:"embedded;embeddedPrefix:build_"`
	Run    DockerTask `gorm:"embedded;embeddedPrefix:run_"`
	Status TaskStatus `gorm:"embedded"`
}

type MatchDetail struct {
	Template string
}

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

func GetGamePrivilege(gameId uint, userId uint) (GamePrivilege, error) {
	// TODO: implement
	return GamePrivilegeRegistered, nil
}
