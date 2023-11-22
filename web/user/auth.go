package user

import (
	"golang.org/x/crypto/bcrypt"

	"hiper-backend/model"
)

func CheckPassword(username string, password string) bool {
	user, err := model.GetUserByUsername(username, "Password")
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return err == nil
}
