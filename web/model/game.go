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

// CRUD: Create

func (g *Game) AfterCreate(tx *gorm.DB) (err error) {
	g.GameId = g.ID
	return tx.Save(g).Error
}

func CreateGame(game *Game, adminIDs []uint) error {
	admins := make([]User, len(adminIDs))
	for i, id := range adminIDs {
		admins[i] = User{Model: gorm.Model{ID: id}}
	}
	game.Admins = admins
	return db.Create(game).Error
}

// CRUD: Read

// Here readme is truncated to 100 characters.
func GetGames(fields ...string) ([]Game, error) {
	var games []Game
	err := db.Select(fields).Find(&games).Error
	return games, err
}

func GetGameById(id uint, fields ...string) (Game, error) {
	var game Game
	err := db.Select(fields).First(&game, id).Error
	return game, err
}

// CRUD: Update

func UpdateGameById(id uint, updates map[string]interface{}) error {
	return db.Model(&Game{}).Where("id = ?", id).Updates(updates).Error
}

// CRUD: Delete

func DeleteGameById(id uint) error {
	return db.Delete(&Game{}, id).Error
}

// associations

// Note: Game doesn't need registration
func (g *Game) GetPrivilege(userId uint) (GamePrivilege, error) {
	var count int64
	err := db.Table("game_admins").Where("game_id = ? AND user_id = ?", g.ID, userId).Count(&count).Error
	if err != nil {
		return "", err
	}
	if count > 0 {
		return GamePrivilegeAdmin, nil
	}
	return GamePrivilegeRegistered, nil
}

// admin

func (g *Game) AddAdmin(userId uint) error {
	user := User{Model: gorm.Model{ID: userId}}
	return db.Model(g).Association("Admins").Append(&user)
}

func (g *Game) GetAdmins() ([]User, error) {
	var admins []User
	err := db.Model(g).Association("Admins").Find(&admins)
	return admins, err
}

func (g *Game) RemoveAdmin(userId uint) error {
	user := User{Model: gorm.Model{ID: userId}}
	return db.Model(g).Association("Admins").Delete(&user)
}

// contestant

// GetContestants returns all contestants in the game.
// By default, sorted by points in descending order.
func (g *Game) GetContestants(fields ...string) ([]Contestant, error) {
	return getContestants(map[string]interface{}{"game_id": g.ID, "contest_id": 0}, fields...)
}

// ai

func (g *Game) GetAis(query QueryParams, preload bool) (ais []Ai, count int64, err error) {
	if query.Filter == nil {
		query.Filter = make(map[string]interface{})
	}
	query.Filter["game_id"] = g.ID
	query.Filter["contest_id"] = 0
	return GetAis(query, preload)
}

func (g *Game) GetAiById(id uint, preload bool, fields ...string) (Ai, error) {
	return getAi(map[string]interface{}{"game_id": g.ID, "contest_id": 0, "Number": id}, preload, fields...)
}

// match

func (g *Game) GetMatches(query QueryParams) (matches []Match, count int64, err error) {
	if query.Filter == nil {
		query.Filter = make(map[string]interface{})
	}
	query.Filter["game_id"] = g.ID
	query.Filter["contest_id"] = 0
	return GetMatches(query)
}

// sdk

func (g *Game) GetSdks(fields ...string) ([]Sdk, error) {
	return GetSdks(map[string]interface{}{"game_id": g.ID, "contest_id": 0}, fields...)
}
