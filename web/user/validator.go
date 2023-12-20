package user

import (
	"hiper-backend/model"
	"net/url"
	"regexp"
)

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
