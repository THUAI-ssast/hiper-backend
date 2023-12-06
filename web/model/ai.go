package model

import "gorm.io/gorm"

type Ai struct {
	gorm.Model
	GameId uint `gorm:"index"`
	// ContestId is 0 if the AI is in a game instead of any contest.
	ContestId uint `gorm:"index"`
	// Number is a unique identifier for each AI within a game or contest.
	Number uint `gorm:"index"`

	UserId uint `gorm:"index"`
	User   User `gorm:"foreignKey:UserId"`
	SdkId  uint
	Sdk    Sdk `gorm:"foreignKey:SdkId"`

	Note   string
	Status TaskStatus `gorm:"embedded;embeddedPrefix:task_"`
}

// TODO: add CRUD functions for ai
