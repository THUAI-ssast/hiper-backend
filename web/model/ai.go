package model

import (
	"errors"
	"path/filepath"
	"strconv"
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

// If preload is true, the following fields will be preloaded:
// Sdk.ID, Sdk.Name
// User.AvatarURL, User.Nickname, User.Username
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

func GetAiByID(id uint, preload bool) (ai Ai, err error) {
	tx := db.Where("id = ?", id)
	if preload {
		tx = addPreloadsForAi(tx)
	}
	err = tx.First(&ai).Error
	return ai, err


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
	if a.ContestId != 0 {
		relativePathWithoutExt = filepath.Join("contest", strconv.Itoa(int(a.ContestId)))
	} else {
		relativePathWithoutExt = filepath.Join("game", strconv.Itoa(int(a.GameId)))
	}
	relativePathWithoutExt = filepath.Join(relativePathWithoutExt, "ais", strconv.Itoa(int(a.Number)), "src")
	return relativePathWithoutExt
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
