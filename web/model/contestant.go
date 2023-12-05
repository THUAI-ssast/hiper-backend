package model

import "gorm.io/gorm"

// TODO

type Contestant struct {
	gorm.Model

	GameId uint `gorm:"uniqueIndex:idx_contestant,priority:1"`
	// ContestId is 0 if the contestant is in a game instead of any contest.
	ContestId uint `gorm:"uniqueIndex:idx_contestant,priority:2"`
	UserId    uint `gorm:"uniqueIndex:idx_contestant,priority:3"`

	Performance string // editable by contest script
	Permissions ContestantPermissions
	Points      int

	AssignedAiId uint
}

type ContestantPermissions struct {
}

func CreateContestant(contestant Contestant) error {
	return db.Create(&contestant).Error
}
