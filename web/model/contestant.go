package model

import "gorm.io/gorm"

type Contestant struct {
	gorm.Model
	GameId uint `gorm:"index"`
	// ContestId is 0 if the contestant is in a game instead of any contest.
	ContestId uint    `gorm:"index"`
	Contest   Contest `gorm:"foreignKey:ContestId"`
	UserId    uint    `gorm:"index"`

	Performance string                // editable by contest script
	Permissions ContestantPermissions `gorm:"embedded"`
	Points      int
	// AssignedAiId is 0 if the contestant doesn't assign an AI.
	AssignedAiId uint
}

func (c *Contestant) BeforeCreate(tx *gorm.DB) (err error) {
	// fill GameId from ContestId
	if c.ContestId != 0 && c.GameId == 0 {
		var gameId uint
		if err = tx.Model(&Contest{}).Select("game_id").First(&Contest{}, c.ContestId).Error; err != nil {
			return err
		}
		c.GameId = gameId
	}
	return nil
}

type ContestantPermissions struct {
	AssignAiEnabled    bool `gorm:"default:true"`
	PublicMatchEnabled bool `gorm:"default:true"`
}

// CRUD: Create

// CreateContestant creates a contestant.
// Either GameId or ContestId must be filled.
// UserId must be filled.
func CreateContestant(contestant *Contestant) error {
	return db.Create(contestant).Error
}

// CRUD: Read

// Sorted by points in descending order.
func getContestants(filter map[string]interface{}, preload preloadQuery) ([]Contestant, error) {
	var contestants []Contestant
	db := db.Where(filter)
	if preload.Table != "" {
		db = db.Preload(preload.Table, func(db *gorm.DB) *gorm.DB {
			return db.Select(preload.Columns)
		})
	}
	err := db.Order("points desc").Find(&contestants).Error
	return contestants, err
}

func getContestant(condition map[string]interface{}) (Contestant, error) {
	var contestant Contestant
	err := db.Where(condition).First(&contestant).Error
	return contestant, err
}

// CRUD: Update

func updateContestant(condition map[string]interface{}, updates map[string]interface{}) error {
	return db.Model(&Contestant{}).Where(condition).Updates(updates).Error
}

func (c *Contestant) Update(updates map[string]interface{}) error {
	return db.Model(c).Updates(updates).Error
}

// CRUD: Delete

func deleteContestant(condition map[string]interface{}) error {
	return db.Delete(&Contestant{}, condition).Error
}

func (c *Contestant) Delete() error {
	return db.Delete(c).Error
}
