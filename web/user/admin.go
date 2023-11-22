package user

import (
	"hiper-backend/model"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// UpsertSuperAdmin upserts the super admin user.
// Super admin must be the first user and its username must be "admin".
// The function gets info from config.
// password required.
func UpsertSuperAdmin() {
	password := viper.GetString("superadmin.password")
	if password == "" {
		panic("superadmin.password is required")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	model.UpsertUser(model.User{
		Username: "admin",
		Password: hashedPassword,
	})
}
