package model

import (
	"errors"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type Ai struct {
	gorm.Model
	GameID uint `gorm:"index"`
	// ContestID is 0 if the AI is in a game instead of any contest.
	ContestID uint `gorm:"index"`
	// Number is a unique identifier for each AI within a game or contest.
	Number uint `gorm:"index"`

	UserID uint `gorm:"index"`
	User   User `gorm:"foreignKey:UserID"`
	SdkID  uint
	Sdk    Sdk `gorm:"foreignKey:SdkID"`

	Note   string
	Status TaskStatus `gorm:"embedded;embeddedPrefix:task_"`
}

func (a *Ai) BeforeCreate(tx *gorm.DB) (err error) {
	// Fill GameID from ContestID
	if a.ContestID != 0 && a.GameID == 0 {
		var gameID uint
		if err = tx.Model(&Contest{}).Select("game_id").First(&Contest{}, a.ContestID).Error; err != nil {
			return err
		}
		a.GameID = gameID
	}
	// Fill Number
	var maxNumber uint
	if err = tx.Model(&Ai{}).Where("game_id = ? AND contest_id = ?", a.GameID, a.ContestID).Pluck("MAX(number)", &maxNumber).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		maxNumber = 0 // If no rows are found, set maxNumber to 0
	}
	a.Number = maxNumber + 1
	return nil
}

// CRUD: Create

func createAi(ai *Ai) error {
	return db.Create(ai).Error
}

// CRUD: Read

// If preload is true, sdk and user will be preloaded, but only some basic fields.
// sdk: id, name
// user: avatar_url, nickname, username
func GetAis(query QueryParams, preload bool) (ais []Ai, count int64, err error) {
	db := db.Select(query.Fields).Where(query.Filter)
	if preload {
		db = db.Preload("Sdk", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("avatar_url", "nickname", "username")
		})
	}
	err = db.Limit(query.Limit).Offset(query.Offset).Find(&ais).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return ais, count, nil
}

// If preload is true, sdk and user will be preloaded, but only some basic fields.
// sdk: id, name
// user: avatar_url, nickname, username
func getAi(condition map[string]interface{}, preload bool, fields ...string) (Ai, error) {
	var ai Ai
	db := db.Select(fields).Where(condition)
	if preload {
		db = db.Preload("Sdk", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("avatar_url", "nickname", "username")
		})
	}
	err := db.First(&ai).Error
	return ai, err
}

// CRUD: Update

func updateAi(condition map[string]interface{}, updates map[string]interface{}) error {
	return db.Model(&Ai{}).Where(condition).Updates(updates).Error
}

// associations

// ai file

// SaveFile saves the AI file to the file system, renamed to src, with the same extension.
// If there is no extension, please pass "".
func (a *Ai) SaveFile(content []byte, ext string) error {
	relativePath := a.getRelativePathWithoutExt()
	if ext != "" {
		relativePath += "." + ext
	}
	return saveFile(relativePath, content)
}

// GetFile returns the AI file from the file system.
func (a *Ai) GetFile() ([]byte, error) {
	relativePathWithoutExt := a.getRelativePathWithoutExt()
	return getFileWithAutoExt(relativePathWithoutExt)
}

func (a *Ai) getRelativePathWithoutExt() string {
	var relativePathWithoutExt string
	// Determine whether the AI is in a contest or a game
	if a.ContestID != 0 {
		relativePathWithoutExt = filepath.Join("contest", strconv.Itoa(int(a.ContestID)))
	} else {
		relativePathWithoutExt = filepath.Join("game", strconv.Itoa(int(a.GameID)))
	}
	relativePathWithoutExt = filepath.Join(relativePathWithoutExt, "ais", strconv.Itoa(int(a.Number)), "src")
	return relativePathWithoutExt
}
