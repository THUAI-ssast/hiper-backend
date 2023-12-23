package model

import (
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	BaseContest BaseContest `gorm:"foreignKey:ID"`

	Metadata Metadata `gorm:"embedded"`
	Admins   []User   `gorm:"many2many:game_admins;"`

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

func (g *Game) Create(adminIDs []uint) error {
	// link a base contest or create a new one
	if g.ID != 0 {
		if err := db.First(&g.BaseContest, g.ID).Error; err != nil {
			return err
		}
	} else {
		if err := db.Create(&g.BaseContest).Error; err != nil {
			return err
		}
	}
	// create game
	g.ID = g.BaseContest.ID
	for _, id := range adminIDs {
		user := User{Model: gorm.Model{ID: id}}
		g.Admins = append(g.Admins, user)
	}
	if err := db.Create(g).Error; err != nil {
		return err
	}
	// update base contest's game_id
	return db.Model(&g.BaseContest).Update("game_id", g.ID).Error
}

// CRUD: Read

func GetGames(fields ...string) (games []Game, err error) {
	tx := db.Preload("BaseContest", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "game_id", "states")
	})
	if len(fields) > 0 {
		tx = tx.Select(fields)
	}
	err = tx.Find(&games).Error
	return
}

func GetGameByID(id uint, fields ...string) (game Game, err error) {
	err = db.Preload("BaseContest").First(&game, id).Error
	return
}

// CRUD: Update

func UpdateGameByID(id uint, updates map[string]interface{}) error {
	return db.Model(&Game{}).Where("id = ?", id).Updates(updates).Error
}

func (g *Game) Update(updates map[string]interface{}) error {
	return db.Model(g).Updates(updates).Error
}

// CRUD: Delete

func DeleteGameByID(id uint) error {
	return db.Delete(&Game{}, id).Error
}

func (g *Game) Delete() error {
	return db.Delete(g).Error
}

// associations

// Note: Game doesn't need registration
func (g *Game) GetPrivilege(userID uint) (GamePrivilege, error) {
	var count int64
	err := db.Table("game_admins").Where("game_id = ? AND user_id = ?", g.ID, userID).Count(&count).Error
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
	user := User{Model: gorm.Model{ID: userId}, Password: []byte{1}}
	return db.Model(g).Association("Admins").Append(&user)
}

func (g *Game) GetAdmins() ([]User, error) {
	var admins []User
	err := db.Model(g).Association("Admins").Find(&admins)
	return admins, err
}

func (g *Game) RemoveAdmin(userId uint) error {
	user := User{Model: gorm.Model{ID: userId}, Password: []byte{1}}
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
