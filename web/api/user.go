package api

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"hiper-backend/config"
	"hiper-backend/mail"

	_ "github.com/go-sql-driver/mysql"
)

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
	timeUnix := time.Now().Unix()
	sqlStr := fmt.Sprintf(`insert into codegiven(email,code,time) values("%s","%s",%d)`, email, code, timeUnix) //sql语句
	ret, err := config.Db.Exec(sqlStr)                                                                          //执行sql语句
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
	mail.MailTo(email, code)
}

func code_match(code string, email string) bool {
	//1.查询记录的sql语句
	sqlStr := fmt.Sprintf("select * from codegiven where email=%s;", email)
	//2.执行
	rst, ok := config.Query(sqlStr) //从连接池里取一个连接出来去数据库查询记录
	if ok {
		timeUnix := time.Now().Unix()
		timeThen, _ := strconv.ParseInt(rst[0]["time"], 10, 64)
		return code == rst[0]["code"] && timeUnix-timeThen < 300
	} else {
		return false
	}
}

func set_user(email string, password string) int {
	//删除数据库db中的数据
	sqlStr := `delete from codegiven where email=?`
	_, err := config.Db.Exec(sqlStr, email)
	if err != nil {
		return -1
	}
	return 1 //TODO: set user in database
}

func verify_password(password string) bool {
	expr := `^[0-9A-Za-z!@#$%^&*]{8,16}$`
	reg := regexp.MustCompile(expr)
	m := reg.MatchString(password)
	return m
}

func email_not_exist(email string) bool {
	//1.查询记录的sql语句
	sqlStr := fmt.Sprintf("select * from user where email=%s;", email)
	//2.执行
	_, ok := config.Query(sqlStr) //从连接池里取一个连接出来去数据库查询记录
	return !ok
}
