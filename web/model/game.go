package model

import (
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	BaseContest BaseContest `gorm:"foreignKey:ID"`

	Metadata Metadata `gorm:"embedded"`

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
		if err := g.BaseContest.Create(adminIDs); err != nil {
			return err
		}
	}
	// create game
	g.ID = g.BaseContest.ID
	if err := db.Create(g).Error; err != nil {
		return err
	}
	// update base contest's game_id
	return db.Model(&g.BaseContest).Update("game_id", g.ID).Error
}

// CRUD: Read

func GetGames(fields ...string) (games []Game, err error) {
	tx := db.Preload("BaseContest", func(db *gorm.DB) *gorm.DB {
		return db.Select(baseContestBaseFields)
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
	err := db.Table("base_contest_admins").Where("base_contest_id = ? AND user_id = ?", g.ID, userID).Count(&count).Error
	if err != nil {
		return "", err
	}
	if count > 0 {
		return GamePrivilegeAdmin, nil
	}
	return GamePrivilegeRegistered, nil
}

// admin

func (g *Game) AddAdmin(userID uint) error {
	user := User{Model: gorm.Model{ID: userID}}
	return db.Model(g).Association("Admins").Append(&user)
}

func (g *Game) GetAdmins() ([]User, error) {
	var admins []User
	err := db.Model(g).Association("Admins").Find(&admins)
	return admins, err
}

func (g *Game) RemoveAdmin(userID uint) error {
	user := User{Model: gorm.Model{ID: userID}}
	return db.Model(g).Association("Admins").Delete(&user)
}

// TODO:game logic files
