package model

import (
	"gorm.io/gorm"
)

type Contest struct {
	gorm.Model
	BaseContest

	Admins []User `gorm:"many2many:contest_admins;"`

	// TODO
	// registration
}

type ContestPrivilege string

const (
	ContestPrivilegeAdmin        ContestPrivilege = "admin"
	ContestPrivilegeRegistered   ContestPrivilege = "registered"
	ContestPrivilegeUnregistered ContestPrivilege = "unregistered"
)

// TODO: add CRUD functions for contest
