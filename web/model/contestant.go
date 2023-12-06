package model

import "gorm.io/gorm"

type Contestant struct {
	gorm.Model
	GameId uint `gorm:"index"`
	// ContestId is 0 if the contestant is in a game instead of any contest.
	ContestId uint `gorm:"index"`
	UserId    uint `gorm:"index"`

	Performance string                // editable by contest script
	Permissions ContestantPermissions `gorm:"embedded"`
	Points      int
	// AssignedAiId is 0 if the contestant doesn't assign an AI.
	AssignedAiId uint
}

func (c *Contestant) BeforeCreate(tx *gorm.DB) (err error) {
	// fill GameId from ContestId
	if c.ContestId != 0 && c.GameId == 0 {
		var gameId uint
		if err = tx.Model(&Contest{}).Select("game_id").First(&Contest{}, c.ContestId).Error; err != nil {
			return err
		}
		c.GameId = gameId
	}
	return nil
}

type ContestantPermissions struct {
	AssignAiEnabled    bool `gorm:"default:true"`
	PublicMatchEnabled bool `gorm:"default:true"`
}

// CreateContestant creates a contestant.
// Either GameId or ContestId must be filled.
func CreateContestant(contestant Contestant) error {
	return db.Create(&contestant).Error
}

// TODO: add CRUD functions for contestant
