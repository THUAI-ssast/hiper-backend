package model

import (
	"gorm.io/gorm"
)

type Ai struct {
	gorm.Model
	BaseContestID uint `gorm:"index"`

	UserID uint `gorm:"index"`
	User   User `gorm:"foreignKey:UserID"`
	SdkID  uint
	Sdk    Sdk `gorm:"foreignKey:SdkID"`

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

func (a *Ai) Create() error {
	return db.Create(a).Error
}

// CRUD: Read

// If preload is true, sdk and user will be preloaded, but only some basic fields.
// sdk: id, name
// user: avatar_url, nickname, username
func GetAis(query QueryParams, preload bool) (ais []Ai, count int64, err error) {
	tx := db.Select(query.Fields).Where(query.Filter)
	if preload {
		tx = addPreloadsForAi(tx)
	}
	tx = tx.Session(&gorm.Session{})

	if query.Limit == 0 {
		query.Limit = 20
	}
	if err = tx.Limit(query.Limit).Offset(query.Offset).Find(&ais).Error; err != nil {
		return nil, 0, err
	}
	if err = tx.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return ais, count, nil
}

func GetAiByID(id uint, preload bool) (ai Ai, err error) {
	tx := db.Where("id = ?", id)
	if preload {
		tx = addPreloadsForAi(tx)
	}
	err = tx.First(&ai).Error
	return ai, err
}

func addPreloadsForAi(tx *gorm.DB) *gorm.DB {
	return tx.Preload("Sdk", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("avatar_url", "nickname", "username")
	})
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
