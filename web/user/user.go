package user

import (
	"fmt"
	"hiper-backend/mail"
	"hiper-backend/model"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
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

func IsCodeMatch(code string, email string) bool {
	storedCode, err := model.GetVerificationCode(email)
	if err != nil {
		return false
	}
	return code == storedCode
}

var passwordRegexp = regexp.MustCompile(`^[0-9A-Za-z!@#$%^&*()\[\]{}<>.,;:?/|~]{8,16}$`)

func IsValidPassword(password string) bool {
	return passwordRegexp.MatchString(password)
}

func IsValidURL(urlStr string) bool {
	if urlStr == "" {
		return true
	}
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
