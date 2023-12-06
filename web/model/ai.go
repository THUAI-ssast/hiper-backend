package model

import "gorm.io/gorm"

type Ai struct {
	gorm.Model
	GameId uint `gorm:"uniqueIndex:idx_ai,priority:1"`
	// ContestId is 0 if the AI is in a game instead of any contest.
	ContestId uint `gorm:"uniqueIndex:idx_ai,priority:2"`
	// Number is a unique identifier for each AI within a game or contest.
	Number uint `gorm:"uniqueIndex:idx_ai,priority:3"`

	UserId uint `gorm:"index"`
	User   User `gorm:"foreignKey:UserId"`
	SdkId  uint
	Sdk    Sdk `gorm:"foreignKey:SdkId"`

	Note   string
	Status TaskStatus `gorm:"embedded;embeddedPrefix:task_"`
}

// TODO: add CRUD functions for ai
