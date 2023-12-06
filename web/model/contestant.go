package model

import "gorm.io/gorm"

type Contestant struct {
	gorm.Model
	GameId uint `gorm:"uniqueIndex:idx_contestant,priority:1"`
	// ContestId is 0 if the contestant is in a game instead of any contest.
	ContestId uint `gorm:"uniqueIndex:idx_contestant,priority:2"`
	UserId    uint `gorm:"uniqueIndex:idx_contestant,priority:3"`

	Performance string                // editable by contest script
	Permissions ContestantPermissions `gorm:"embedded"`
	Points      int
	// AssignedAiId is 0 if the contestant doesn't assign an AI.
	AssignedAiId uint
}

type ContestantPermissions struct {
	AssignAiEnabled    bool `gorm:"default:true"`
	PublicMatchEnabled bool `gorm:"default:true"`
}

func CreateContestant(contestant Contestant) error {
	return db.Create(&contestant).Error
}

// TODO: add CRUD functions for contestant
