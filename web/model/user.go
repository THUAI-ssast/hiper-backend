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

	GameAdmins    []Game    `gorm:"many2many:game_admins;"`
	ContestAdmins []Contest `gorm:"many2many:contest_admins;"`
}

// CRUD: Create

// CreateUser creates a user. `user`'s ID will be updated if the operation succeeds.
func CreateUser(user *User) error {
	return db.Create(user).Error
}

// CRUD: Read

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

func UpdateUserById(id uint, updates map[string]interface{}) error {
	return updateUser(map[string]interface{}{"id": id}, updates)
}

func UpdateUserByEmail(email string, updates map[string]interface{}) error {
	return updateUser(map[string]interface{}{"email": email}, updates)
}

func updateUser(condition map[string]interface{}, updates map[string]interface{}) error {
	return db.Model(&User{}).Where(condition).Updates(updates).Error
}

// associations

// admin

func (u *User) GetGameAdmins(fields ...string) ([]Game, error) {
	var games []Game
	err := db.Model(u).Select(fields).Association("GameAdmins").Find(&games)
	return games, err
}

func (u *User) GetContestAdmins(fields ...string) ([]Contest, error) {
	var contests []Contest
	err := db.Model(u).Select(fields).Association("ContestAdmins").Find(&contests)
	return contests, err
}

// contestant

// GetContestants returns the contestants of the user.
// If `preload` is true, contest information will be preloaded.
func (u *User) GetContestants(preload bool) ([]Contestant, error) {
	if preload {
		return getContestants(map[string]interface{}{"user_id": u.ID}, preloadQuery{
			"Contest", []string{"id", "cover_url", "title", "public_match_enabled"},
		})
	} else {
		return getContestants(map[string]interface{}{"user_id": u.ID}, preloadQuery{})
	}
}

// Irregular CRUD

// UpsertUser upserts a user.
// If the user exists(determinated by username), update its password.
// If the user does not exist, create it.
func UpsertUser(user User) {
	// Find the user by username or create a new one
	result := db.Where(User{Username: user.Username}).FirstOrCreate(&user)

	// If the user was found, update the password
	if result.RowsAffected > 0 {
		db.Model(&user).Update("password", user.Password)
	}
}
