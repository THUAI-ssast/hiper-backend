package user

import (
	"github.com/spf13/viper"

	"github.com/THUAI-ssast/hiper-backend/model"
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
	hashedPassword := HashPassword(password)

	model.UpsertUser(model.User{
		Username: "admin",
		Email:    "admin@mails.tsinghua.edu.cn",
		Permissions: model.GlobalPermissions{
			CanCreateGameOrContest: true,
		},
		Password: hashedPassword,
	})
}
