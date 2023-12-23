package model

import (
	"errors"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model
	GameID uint `gorm:"index"`
	// ContestID is 0 if the match is in a game instead of any contest.
	ContestID uint `gorm:"index"`
	// Number is a unique identifier for each match within a game or contest.
	Number uint `gorm:"index"`

	Players []Ai `gorm:"many2many:match_ais;"`
	State   TaskState
	Tag     string

	Scores []int `gorm:"serializer:json"`
}

// CRUD: Create

func (m *Match) BeforeCreate(tx *gorm.DB) (err error) {
	// Fill GameID from ContestID
	if m.ContestID != 0 && m.GameID == 0 {
		var gameID uint
		if err = tx.Model(&Contest{}).Select("game_id").First(&Contest{}, m.ContestID).Error; err != nil {
			return err
		}
		m.GameID = gameID
	}
	// Fill Number
	var maxNumber uint
	if err = tx.Model(&Match{}).Where("game_id = ? AND contest_id = ?", m.GameID, m.ContestID).Pluck("MAX(number)", &maxNumber).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		maxNumber = 0 // If no rows are found, set maxNumber to 0
	}
	m.Number = maxNumber + 1
	return nil
}

// TODO
func CreateMatch(match *Match, playerIDs []uint) error {
	players := make([]Ai, len(playerIDs))
	for i, id := range playerIDs {
		players[i] = Ai{Model: gorm.Model{ID: id}}
	}
	match.Players = players
	return db.Create(match).Error
}

// CRUD: Read

func GetMatches(query QueryParams) ([]Match, int64, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}
	var matches []Match
	var count int64
	db := db.Where(query.Filter)
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Offset(query.Offset).Limit(query.Limit).Find(&matches).Error
	return matches, count, err
}

func getMatch(condition map[string]interface{}, fields ...string) (Match, error) {
	var match Match
	err := db.Select(fields).Where(condition).First(&match).Error
	return match, err
}

// associations

// GetLogs returns logs of each player in the match
func (m *Match) GetLogs() ([]string, error) {
	// TODO: implement
	// read logs from file
	return nil, nil
}

// GetReplay returns replay of the match
func (m *Match) GetReplay() (string, error) {
	// TODO: implement
	// read replay from file
	return "", nil
}
