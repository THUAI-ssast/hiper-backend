package model

import "gorm.io/gorm"

type Match struct {
	gorm.Model
	GameId uint `gorm:"index"`
	// ContestId is 0 if the match is in a game instead of any contest.
	ContestId uint `gorm:"index"`
	// Number is a unique identifier for each match within a game or contest.
	Number uint `gorm:"index"`

	Players []Ai `gorm:"many2many:match_ais;"`
	State   TaskState
	Tag     string

	Logs   []string `gorm:"serializer:json"`
	Replay string
	Scores []int `gorm:"serializer:json"`
}

// TODO: add CRUD functions for match
