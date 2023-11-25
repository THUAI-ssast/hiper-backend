package model

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Permissions struct {
	CanCreateGameOrContest bool `gorm:"default: false"`
}

type User struct {
	gorm.Model
	AvatarURL          string
	Bio                string
	ContestsRegistered []Contest `gorm:"many2many:contest_registrations;"`
	Department         string
	Email              string `gorm:"uniqueIndex,not null"`
	Name               string
	Nickname           string      `gorm:"index"`
	Password           []byte      `gorm:"not null"`
	Permissions        Permissions `gorm:"embedded"`
	School             string
	Username           string `gorm:"uniqueIndex,not null"`
}

func CreateUser(user User) error {
	return db.Create(&user).Error
}

func GetUserByUsername(username string, fields ...string) (User, error) {
	return getUser(map[string]interface{}{"username": username}, fields...)
}

func GetUserByEmail(email string, fields ...string) (User, error) {
	return getUser(map[string]interface{}{"email": email}, fields...)
}

func GetUserById(id uint, fields ...string) (User, error) {
	return getUser(map[string]interface{}{"id": id}, fields...)
}

func SearchUsers(keyword string, searchFields []string, resultFields ...string) ([]User, error) {
	var users []User
	query := strings.Join(searchFields, " LIKE ? OR ") + " LIKE ?"
	args := make([]interface{}, len(searchFields))
	for i := range args {
		args[i] = "%" + keyword + "%"
	}
	err := db.Select(resultFields).Where(query, args...).Find(&users).Error
	return users, err
}

func getUser(condition map[string]interface{}, fields ...string) (User, error) {
	var user User
	err := db.Select(fields).Where(condition).First(&user).Error
	return user, err
}

func UpdateUserByUsername(username string, updates map[string]interface{}) error {
	return updateUser(map[string]interface{}{"username": username}, updates)
}

func UpdateUserById(id uint, updates map[string]interface{}) error {
	return updateUser(map[string]interface{}{"id": id}, updates)
}

func UpdateUserByEmail(email string, updates map[string]interface{}) error {
	return updateUser(map[string]interface{}{"email": email}, updates)
}

func updateUser(condition map[string]interface{}, updates map[string]interface{}) error {
	return db.Model(&User{}).Where(condition).Updates(updates).Error
}

// UpsertUser upserts a user.
// If the user exists, update its password.
// If the user does not exist, create it.
func UpsertUser(user User) {
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},
		DoUpdates: clause.AssignmentColumns([]string{"password"}),
	}).Create(&user)
}
