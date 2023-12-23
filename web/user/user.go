package user

import (
	"fmt"
	"math/rand"
	"strings"

	"hiper-backend/mail"
	"hiper-backend/model"
)

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendVerificationCode(email string) error {
	code := GenValidateCode(6)
	if err := model.SaveVerificationCode(code, email, 5); err != nil {
		return err
	}
	if err := mail.SendVerificationCode(code, email); err != nil {
		return err
	}
	return nil
}

// RegisterUser registers a user.
// It returns the user id and an error.
// The error is nil if the registration is successful.
func RegisterUser(username string, email string, password string) (uint, error) {
	hashedPassword := HashPassword(password)
	user := model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	if err := user.Create(); err != nil {
		return 0, err
	}
	return user.ID, nil
}
