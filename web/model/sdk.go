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

func UpdateSdkByID(id uint, updates map[string]interface{}) error {
	return db.Model(&Sdk{}).Where("id = ?", id).Updates(updates).Error
}

func (s *Sdk) Update(updates map[string]interface{}) error {
	return db.Model(s).Updates(updates).Error
}

func getSdk(condition map[string]interface{}, fields ...string) (Sdk, error) {
	var sdk Sdk
	db := db.Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	err := db.First(&sdk).Error
	return sdk, err
}

func GetSdkByID(id uint, fields ...string) (Sdk, error) {
	return getSdk(map[string]interface{}{"id": id}, fields...)
}

// CRUD: Delete

func DeleteSdkByID(id uint) error {
	return db.Delete(&Sdk{}, id).Error
}

func (s *Sdk) Delete() error {
	return db.Delete(s).Error
}

// sdk file
