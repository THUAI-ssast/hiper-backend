package model

import (
	"strings"

	"gorm.io/gorm"
)

type GlobalPermissions struct {
	CanCreateGameOrContest bool `gorm:"default: false"`
}

type User struct {
	gorm.Model
	AvatarURL   string
	Bio         string
	Department  string
	Email       string `gorm:"uniqueIndex,not null"`
	Name        string
	Nickname    string            `gorm:"index"`
	Password    []byte            `gorm:"not null"`
	Permissions GlobalPermissions `gorm:"embedded"`
	School      string
	Username    string `gorm:"uniqueIndex,not null"`

	ContestRegistered []Contest `gorm:"many2many:contest_registrations;"`
}

// CRUD: Create

// Necessary fields: Username, Email, Password
// Optional fields: AvatarURL, Bio, Department, Name, Nickname, School, Permissions
func (u *User) Create() error {
	return db.Create(u).Error
}

// CRUD: Read

func GetUserByUsername(username string, fields ...string) (User, error) {
	return getUser(map[string]interface{}{"username": username}, fields...)
}

func GetUserByEmail(email string, fields ...string) (User, error) {
	return getUser(map[string]interface{}{"email": email}, fields...)
}

func GetUserByID(id uint, fields ...string) (User, error) {
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
	db := db.Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	err := db.First(&user).Error
	return user, err
}

// CRUD: Update

func UpdateUserByUsername(username string, updates map[string]interface{}) error {
	return updateUser(map[string]interface{}{"username": username}, updates)
}

func UpdateUserByID(id uint, updates map[string]interface{}) error {
	return updateUser(map[string]interface{}{"id": id}, updates)
}

func UpdateUserByEmail(email string, updates map[string]interface{}) error {
	return updateUser(map[string]interface{}{"email": email}, updates)
}

func updateUser(condition map[string]interface{}, updates map[string]interface{}) error {
	return db.Model(&User{}).Where(condition).Updates(updates).Error
}

func (u *User) Update(updates map[string]interface{}) error {
	return db.Model(u).Updates(updates).Error
}

// associations

func (u *User) GetContestRegistered(fields ...string) ([]Contest, error) {
	var contests []Contest
	err := db.Model(u).Select(fields).Preload("BaseContest", func(db *gorm.DB) *gorm.DB {
		return db.Select(baseContestBaseFields)
	}).Association("ContestRegistered").Find(&contests)
	return contests, err
}

// Irregular CRUD

// UpsertUser upserts a user.
// If the user exists(determined by username), update its password.
// If the user does not exist, create it.
func UpsertUser(user User) {
	// Find the user by username or create a new one
	result := db.Where(User{Username: user.Username}).FirstOrCreate(&user)

	// If the user was found, update the password
	if result.RowsAffected > 0 {
		db.Model(&user).Update("password", user.Password)
	}
}
