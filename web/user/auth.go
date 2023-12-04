package user

import (
	"golang.org/x/crypto/bcrypt"

	"hiper-backend/model"
)

func CheckPasswordByUsername(username string, password string) (uint, bool) {
	user, err := model.GetUserByUsername(username, "Password", "ID")
	if err != nil {
		return 0, false
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return user.ID, err == nil
}

func CheckPasswordByEmail(email string, password string) (uint, bool) {
	user, err := model.GetUserByEmail(email, "Password", "ID")
	if err != nil {
		return 0, false
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return user.ID, err == nil
}

func HashPassword(password string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash
}
