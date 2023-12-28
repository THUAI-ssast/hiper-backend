package model

import (
	"errors"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model
	BaseContestID uint `gorm:"index"`

	Players []Ai `gorm:"many2many:match_ais;"`
	State   TaskState
	Tag     string

	Scores []int `gorm:"serializer:json"`

	// snapshot fields
	Users []User `gorm:"-"` // not stored in database
}

// CRUD: Create

// Necessary fields: BaseContestID
// Optional fields: Tag
func (m *Match) Create(playerIDs []uint) error {
	if len(playerIDs) == 0 {
		return errors.New("no players")
	}

	if err := db.Create(m).Error; err != nil {
		return err
	}

	for _, playerID := range playerIDs {
		if err := db.Exec("INSERT INTO match_ais (match_id, ai_id) VALUES (?, ?)", m.ID, playerID).Error; err != nil {
			return err
		}
	}
	return nil
}

// CRUD: Read

// Sorted by id in descending order.
// If preload is true, the following fields will be preloaded:
// Ai.ID, Ai.Sdk.ID, Ai.Sdk.Name
// User.AvatarURL, User.Nickname, User.Username
func GetMatches(query QueryParams, preload bool) (matches []Match, count int64, err error) {
	tx := db.Order("id DESC")
	if preload {
		tx = addPreloadsForMatch(tx)
	}
	if query.Limit == 0 {
		query.Limit = 20
	}
	count, err = paginate(tx, query, &matches)
	return matches, count, nil
}

func GetMatchByID(id uint, preload bool) (match Match, err error) {
	tx := db.Where("id = ?", id)
	if preload {
		tx = addPreloadsForMatch(tx)
	}
	err = tx.First(&match).Error
	return
}

func addPreloadsForMatch(tx *gorm.DB) *gorm.DB {
	return tx.Preload("Players", func(db *gorm.DB) *gorm.DB {
		return db.Select("id").Preload("Sdk", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("avatar_url", "nickname", "username")
		})
	})
}

// CRUD: Update

func UpdateMatchByID(id uint, updates map[string]interface{}) error {
	return db.Model(&Match{}).Where("id = ?", id).Updates(updates).Error
}

func (m *Match) Update(updates map[string]interface{}) error {
	return db.Model(m).Updates(updates).Error
}

// associations

// match files
