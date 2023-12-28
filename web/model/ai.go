package model

import (
	"gorm.io/gorm"
)

type Ai struct {
	gorm.Model
	BaseContestID uint `gorm:"not null;index"`

	UserID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID"`
	SdkID  uint `gorm:"not null"`
	Sdk    Sdk  `gorm:"foreignKey:SdkID"`

	Note   string
	Status TaskStatus `gorm:"embedded;embeddedPrefix:task_"`

	// snapshot fields
	GameID uint
}

func (a *Ai) BeforeCreate(tx *gorm.DB) (err error) {
	// Fill GameID from BaseContestID
	var bc BaseContest
	if err = tx.Model(&BaseContest{}).Select("game_id").First(&bc, a.BaseContestID).Error; err != nil {
		return err
	}
	a.GameID = bc.GameID
	return nil
}

// CRUD: Create

// Necessary fields: BaseContestID, UserID, SdkID
// Optional fields: Note
func (a *Ai) Create() error {
	return db.Create(a).Error
}

// CRUD: Read

// If preload is true, the following fields will be preloaded:
// Sdk, User
func GetAis(query QueryParams, preload bool) (ais []Ai, count int64, err error) {
	tx := db.Order("id DESC")
	if preload {
		tx = addPreloadsForAi(tx)
	}
	if query.Limit == 0 {
		query.Limit = 20
	}
	count, err = paginate(tx, query, &ais)
	return ais, count, nil
}

// If preload is true, the following fields will be preloaded:
// Sdk, User
func GetAiByID(id uint, preload bool) (ai Ai, err error) {
	tx := db.Where("id = ?", id)
	if preload {
		tx = addPreloadsForAi(tx)
	}
	err = tx.First(&ai).Error
	return ai, err
}

func addPreloadsForAi(tx *gorm.DB) *gorm.DB {
	return tx.Preload("Sdk").Preload("User")
}

// CRUD: Update

func UpdateAiByID(id uint, updates map[string]interface{}) error {
	return db.Model(&Ai{}).Where("id = ?", id).Updates(updates).Error
}

func (a *Ai) Update(updates map[string]interface{}) error {
	return db.Model(a).Updates(updates).Error
}

// associations

// ai file
