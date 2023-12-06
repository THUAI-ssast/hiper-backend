package model

import "gorm.io/gorm"

type Sdk struct {
	gorm.Model
	GameId uint `gorm:"index"`
	// ContestId is 0 if the SDK is in a game instead of any contest.
	ContestId uint `gorm:"index"`
	// Number is a unique identifier for each SDK within a game or contest.
	Number uint `gorm:"index"`

	Name   string `gorm:"unique"`
	Readme string

	BuildAi DockerTask `gorm:"embedded;embeddedPrefix:build_ai_"`
	RunAi   DockerTask `gorm:"embedded;embeddedPrefix:run_ai_"`
}

// TODO: add CRUD functions for sdk
