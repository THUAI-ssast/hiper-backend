package model

import (
	"gorm.io/gorm"
)

type Contest struct {
	gorm.Model
	BaseContest BaseContest `gorm:"foreignKey:ID"`

	Metadata Metadata `gorm:"embedded"`
	Admins   []User   `gorm:"many2many:contest_admins;"`

	Registration    Registration `gorm:"embedded"`
	RegisteredUsers []User       `gorm:"many2many:contest_registrations;"`
}

type ContestPrivilege string

const (
	ContestPrivilegeAdmin        ContestPrivilege = "admin"
	ContestPrivilegeRegistered   ContestPrivilege = "registered"
	ContestPrivilegeUnregistered ContestPrivilege = "unregistered"
)

type Registration struct {
	RegistrationEnabled bool `gorm:"default: false"`
	Password            string
}

// TODO: add CRUD functions for contest

func (c *Contest) Create(gameID uint, adminIDs []uint) error {
	// link a base contest or create a new one
	if c.ID != 0 {
		if err := db.First(&c.BaseContest, c.ID).Error; err != nil {
			return err
		}
	} else {
		c.BaseContest.GameID = gameID
		if err := db.Create(&c.BaseContest).Error; err != nil {
			return err
		}
	}
	// create contest
	c.ID = c.BaseContest.ID
	for _, id := range adminIDs {
		user := User{Model: gorm.Model{ID: id}}
		c.Admins = append(c.Admins, user)
	}
	if err := db.Create(c).Error; err != nil {
		return err
	}
	return nil
}

// association

func (c *Contest) GetPrivilege(userID uint) (ContestPrivilege, error) {
	// check if the user is an admin
	var count int64
	if err := db.Table("contest_admins").Where("contest_id = ? AND user_id = ?", c.ID, userID).Count(&count).Error; err != nil {
		return "", err
	}
	if count > 0 {
		return ContestPrivilegeAdmin, nil
	}
	// check if the user is registered
	if err := db.Table("contest_registrations").Where("contest_id = ? AND user_id = ?", c.ID, userID).Count(&count).Error; err != nil {
		return "", err
	}
	if count > 0 {
		return ContestPrivilegeRegistered, nil
	}
	return ContestPrivilegeUnregistered, nil
}
