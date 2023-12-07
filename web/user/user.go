package user

import (
	"errors"
	"fmt"
	"hiper-backend/mail"
	"hiper-backend/model"
	"io"
	"math/rand"
	netmail "net/mail"
	"net/url"
	"regexp"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/spf13/viper"
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

func FindVerificationCode() (string, error) {
	// Connect to the server
	host := viper.GetString("mail.host")
	port := viper.GetInt("mail.port")
	server := fmt.Sprintf("%s:%d", host, port)
	c, err := client.DialTLS(server, nil)
	if err != nil {
		return "", err
	}
	defer c.Logout()

	// Login
	if err := c.Login(viper.GetString("mail.username"), viper.GetString("mail.password")); err != nil {
		return "", err
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return "", err
	}

	// Get the last message
	if mbox.Messages == 0 {
		return "", errors.New("No message in mailbox")
	}
	seqset := new(imap.SeqSet)
	seqset.AddNum(mbox.Messages)

	// Get the whole message body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	msg := <-messages
	if msg == nil {
		return "", errors.New("Server didn't returned message")
	}

	r := msg.GetBody(section)
	if r == nil {
		return "", errors.New("Server didn't returned message body")
	}

	// Parse the message body
	mail, err := netmail.ReadMessage(r)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(mail.Body)
	if err != nil {
		return "", err
	}

	// Use regex to find the verification code
	re := regexp.MustCompile(`Your verification code for Hiper is (\d{6}),`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", errors.New("No verification code found")
	}

	return matches[1], nil
}
