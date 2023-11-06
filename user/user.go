package user

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

var codeGiven string

func verify_email(email string) bool {

	return false
}

func send_email(email string) {

}

func Main() {
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := gin.Default()
	r.GET("/users/request-verification-code", func(c *gin.Context) {
		email := c.DefaultQuery("email", "")
		if verify_email(email) {
			send_email(email)
			c.String(200, fmt.Sprintf("verification code sent to %s", email))
		} else {
			c.String(422, fmt.Sprintf("email %s is not valid", email))
		}
	})
	r.GET("/users", func(c *gin.Context) {
		email := c.DefaultQuery("email", "")

		c.String(200, fmt.Sprintf("Hello %s", email))
	})
	r.Run(":8000")
}
