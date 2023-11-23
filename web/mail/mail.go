package mail

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

// SendVerificationCode sends a email with a verification code to the given email address
func SendVerificationCode(code string, email string) error {
	message := `
    	<p> Hello,</p>
	
		<p style="text-indent:2em">Your verification code for Hiper is %s,</p> 
		<p style="text-indent:2em">Please Use it in 5 minutes.</p>
	`

	host := viper.GetString("mail.host")
	port := viper.GetInt("mail.port")
	username := viper.GetString("mail.username")
	password := viper.GetString("mail.password")

	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hiper Verification Code")
	m.SetBody("text/html", fmt.Sprintf(message, code))

	d := gomail.NewDialer(host, port, username, password)
	return d.DialAndSend(m)
}
