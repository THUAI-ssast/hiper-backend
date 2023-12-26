package model

import "gorm.io/gorm"

type Contestant struct {
	gorm.Model
	BaseContestID uint `gorm:"index"`
	UserID        uint `gorm:"index"`
	User          User `gorm:"foreignKey:UserID"`

	Performance string                // editable by contest script
	Permissions ContestantPermissions `gorm:"embedded"`
	Points      int
	// AssignedAiID is 0 if the contestant doesn't assign an AI.
	AssignedAiID uint
	AssignedAi   Ai `gorm:"foreignKey:AssignedAiID"`
}

type ContestantPermissions struct {
	AssignAiEnabled    bool `gorm:"default:true"`
	PublicMatchEnabled bool `gorm:"default:true"`
}

// CRUD: Create

func (c *Contestant) Create() error {
	return db.Create(c).Error
}

// CRUD: Read

// Sorted by points in descending order.
// Currently supported preloads: "User", "AssignedAi"
func GetContestants(filter map[string]interface{}, preloads []PreloadQuery) ([]Contestant, error) {
	var contestants []Contestant
	tx := db.Where(filter).Order("points DESC")
	tx = addPreloads(tx, preloads)
	err := tx.Find(&contestants).Error
	return contestants, err
}

func GetContestant(condition map[string]interface{}, preloads []PreloadQuery) (Contestant, error) {
	contestants, err := GetContestants(condition, preloads)
	if len(contestants) == 0 {
		return Contestant{}, err
	}
	return contestants[0], err
}

func GetContestantByID(id uint, preloads []PreloadQuery) (Contestant, error) {
	var contestant Contestant
	tx := db.Where("id = ?", id)
	tx = addPreloads(tx, preloads)
	err := tx.First(&contestant).Error
	return contestant, err
}

// CRUD: Update

func UpdateContestantByID(id uint, updates map[string]interface{}) error {
	return db.Model(&Contestant{}).Where("id = ?", id).Updates(updates).Error
}

func (c *Contestant) Update(updates map[string]interface{}) error {
	return db.Model(c).Updates(updates).Error
}

// CRUD: Delete

func DeleteContestantByID(id uint) error {
	return db.Delete(&Contestant{}, id).Error
}

func (c *Contestant) Delete() error {
	return db.Delete(c).Error
}
