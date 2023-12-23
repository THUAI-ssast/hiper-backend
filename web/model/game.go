package model

import (
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model

	BaseContestId uint
	BaseContest   BaseContest `gorm:"foreignKey:BaseContestId"`

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

func (g *Game) Create(adminIds []uint) error {
	// link a base contest or create a new one
	if g.BaseContestId != 0 {
		if err := db.First(&g.BaseContest, g.BaseContestId).Error; err != nil {
			return err
		}
	} else {
		if err := db.Create(&g.BaseContest).Error; err != nil {
			return err
		}
	}
	// create game
	g.BaseContestId = g.BaseContest.ID
	for _, id := range adminIds {
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

func GetGames() ([]Game, error) {
	var games []Game
	err := db.Preload("BaseContest", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "game_id", "states")
	}).Find(&games).Error
	return games, err
}

func GetGameById(id uint) (Game, error) {
	var game Game
	err := db.Preload("BaseContest", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "game_id", "states")
	}).First(&game, id).Error
	return game, err
}

// CRUD: Update

func UpdateGameById(id uint, updates map[string]interface{}) error {
	return db.Model(&Game{}).Where("id = ?", id).Updates(updates).Error
}

func (g *Game) Update(updates map[string]interface{}) error {
	return db.Model(g).Updates(updates).Error
}

// CRUD: Delete

func DeleteGameById(id uint) error {
	return db.Delete(&Game{}, id).Error
}

func (g *Game) Delete() error {
	return db.Delete(g).Error
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

// TODO:game logic files
