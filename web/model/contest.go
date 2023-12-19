package model

import (
	"gorm.io/gorm"
)

type Contest struct {
	gorm.Model
	BaseContest

	Admins []User `gorm:"many2many:contest_admins;"`

	Registration Registration `gorm:"embedded"`
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

func GetContestPrivilege(contestId uint, userId uint) (ContestPrivilege, error) {
	// TODO: implement
	return ContestPrivilegeRegistered, nil
}
