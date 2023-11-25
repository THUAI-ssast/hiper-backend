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
	// 生成验证码
	code := GenValidateCode(6)

	// 发送邮件
	if err := mail.SendVerificationCode(code, email); err != nil {
		return err
	}

	// 保存验证码
	if err := model.SaveVerificationCode(code, email, 5); err != nil {
		return err
	}
	return nil
}

func CodeMatch(code string, email string) bool {
	storedCode, err := model.GetVerificationCode(email)
	if err != nil {
		return false
	}
	return code == storedCode
}

func VerifyPassword(password string) bool {
	expr := `^[0-9A-Za-z!@#$%^&*]{8,16}$`
	reg := regexp.MustCompile(expr)
	m := reg.MatchString(password)
	return m
}

func IsValidURL(urlStr string) bool {
	if urlStr == "" {
		return true
	}
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
