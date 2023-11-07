package user

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gomail.v2"
)

var db *sql.DB

type GetRequestVerificationCode struct {
	Email string `json:"email" binding:"required"`
}

type GetRegister struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"verification_code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func verify_email(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func send_email(email string) {
	// 生成随机验证码
	code := GenValidateCode(6)
	sqlStr := fmt.Sprintf(`insert into codegiven(email,code) values("%s",%s)`, email, code) //sql语句
	ret, err := db.Exec(sqlStr)                                                             //执行sql语句
	if err != nil {
		fmt.Printf("insert failed,err:%v\n", err)
		return
	}
	//如果是插入数据的操作，能够拿到插入数据的id
	_, err = ret.LastInsertId()
	if err != nil {
		fmt.Printf("get id failed,err:%v\n", err)
		return
	}
	message := `
    	<p> Hello,</p>
	
		<p style="text-indent:2em">Your verification code for Hiper is %s,</p> 
		<p style="text-indent:2em">Please Use it to log in.</p>
	`

	// QQ 邮箱：
	// SMTP 服务器地址：smtp.qq.com（SSL协议端口：465/994 | 非SSL协议端口：25）
	// 163 邮箱：
	// SMTP 服务器地址：smtp.163.com（端口：25）
	host := "smtp.qq.com"
	port := 25
	userName := "2975587905@qq.com"
	password := "vqbhftsgfpsmdfed"

	m := gomail.NewMessage()
	m.SetHeader("From", userName) // 发件人
	// m.SetHeader("From", "alias"+"<"+userName+">") // 增加发件人别名

	m.SetHeader("To", email)                          // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	m.SetHeader("Subject", "Hiper verification code") // 邮件主题

	// text/html 的意思是将文件的 content-type 设置为 text/html 的形式，浏览器在获取到这种文件时会自动调用html的解析器对文件进行相应的处理。
	// 可以通过 text/html 处理文本格式进行特殊处理，如换行、缩进、加粗等等
	m.SetBody("text/html", fmt.Sprintf(message, code))

	// text/plain的意思是将文件设置为纯文本的形式，浏览器在获取到这种文件时并不会对其进行处理
	// m.SetBody("text/plain", "纯文本")
	// m.Attach("test.sh")   // 附件文件，可以是文件，照片，视频等等
	// m.Attach("lolcatVideo.mp4") // 视频
	// m.Attach("lolcat.jpg") // 照片

	d := gomail.NewDialer(
		host,
		port,
		userName,
		password,
	)
	// 关闭SSL协议认证
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func set_user(email string, password string) int {
	return 1 //TODO: set user in database
}

func password_valid(password string) bool {
	return true //TODO: check password valid
}

func email_not_exist(email string) bool {
	return true //TODO: check email not exist
}

func Main() {
	var err error
	// Init database
	db, err = sql.Open("mysql", "root:lq3525926@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	err = db.Ping() //尝试连接数据库
	if err != nil {
		fmt.Println("connect mysql failed,", err)
		return
	}
	fmt.Println("connect success!")
	//设置数据库连接池的最大连接数
	db.SetMaxIdleConns(10)
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.OpenFile("./user/gin.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := gin.Default()
	r.GET("/users/request-verification-code", func(c *gin.Context) {
		var jsonGetV GetRequestVerificationCode
		if err := c.ShouldBindJSON(&jsonGetV); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := jsonGetV.Email
		if verify_email(email) {
			send_email(email)
			c.JSON(200, gin.H{})
		} else {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "email",
					},
				},
			})
		}
	})
	r.GET("/users", func(c *gin.Context) {
		var jsonGetR GetRegister
		if err := c.ShouldBindJSON(&jsonGetR); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := jsonGetR.Email
		//code := jsonGetR.Code
		password := jsonGetR.Password
		if true { //codeGiven[email] == code {
			if password_valid(password) {
				if email_not_exist(email) {
					userId := set_user(email, password)
					c.JSON(200, gin.H{
						"user_id": userId,
					})
				} else {
					c.JSON(422, gin.H{
						"errors": []gin.H{
							{
								"code":  "already_exists",
								"field": "email",
							},
						},
					})
				}
			} else {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid",
							"field": "password",
						},
					},
				})
			}
		} else {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "verification_code",
					},
				},
			})
		}
	})
	r.Run(":8000")
}
