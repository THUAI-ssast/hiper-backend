package model

import (
	"gorm.io/gorm"
)

type Sdk struct {
	gorm.Model
	BaseContestID uint `gorm:"not null;index"`

	Name   string `gorm:"not null"`
	Readme string

	BuildAi DockerTask `gorm:"embedded;embeddedPrefix:build_ai_"`
	RunAi   DockerTask `gorm:"embedded;embeddedPrefix:run_ai_"`
}

// CRUD: Create

// Necessary fields: BaseContestID, Name
// Optional fields: Readme, BuildAi, RunAi
func (s *Sdk) Create() error {
	return db.Create(s).Error
}

// CRUD: Update

func (s *Sdk) Update(updates map[string]interface{}) error {
	return db.Model(s).Updates(updates).Error
}

// CRUD: Delete

func DeleteSdkByID(id uint) error {
	return db.Delete(&Sdk{}, id).Error
}

func (s *Sdk) Delete() error {
	return db.Delete(s).Error
}

// sdk file
