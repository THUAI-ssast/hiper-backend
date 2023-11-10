package api

import (
	"fmt"
	"hiper-backend/config"
	"hiper-backend/mail"
	"hiper-backend/utils"

	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type GetRequestVerificationCode struct {
	Email string `json:"email" binding:"required"`
}

type GetRegister struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"verification_code" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type GetResetEmail struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"verification_code" binding:"required"`
	NewEmail string `json:"new_email" binding:"required"`
}

type GetResetPassword struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"verification_code" binding:"required"`
	Password string `json:"new_password" binding:"required"`
}

type GetLogin struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

func request_verification_code(c *gin.Context) {
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
}

func register_user(c *gin.Context) {
	var jsonGetR GetRegister
	if err := c.ShouldBindJSON(&jsonGetR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := jsonGetR.Email
	code := jsonGetR.Code
	password := jsonGetR.Password
	username := jsonGetR.Username
	if !code_match(code, email) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "verification_code",
				},
			},
		})
	} else if !verify_password(password) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "password",
				},
			},
		})
	} else if email_exist_codegiven(email) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "already_exists",
					"field": "email",
				},
			},
		})
	} else if username_exist_user(username) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "already_exists",
					"field": "username",
				},
			},
		})
	} else if email_exist_user(email) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "already_exists",
					"field": "email",
				},
			},
		})
	} else {
		userId := set_user(email, password, username)
		if userId == -1 {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":   "missing_field",
						"field":  "email",
						"detail": "Delete from codegiven failed",
					},
				},
			})
		} else {
			c.JSON(200, gin.H{
				"user_id": userId,
			})
		}
	}
}

func reset_email(c *gin.Context) {
	var jsonGetRE GetResetEmail
	if err := c.ShouldBindJSON(&jsonGetRE); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := jsonGetRE.Email
	code := jsonGetRE.Code
	newEmail := jsonGetRE.NewEmail
	if !email_exist_codegiven(email) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "email",
				},
			},
		})
	} else if !verify_email(newEmail) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "new_email",
				},
			},
		})
	} else if !code_match(code, email) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "verification_code",
				},
			},
		})
	} else if email_exist_codegiven(newEmail) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "already_exists",
					"field": "new_email",
				},
			},
		})
	} else {
		update_email(email, newEmail)
		c.JSON(200, gin.H{})
	}
}

func reset_password(c *gin.Context) {
	var jsonGetRP GetResetPassword
	if err := c.ShouldBindJSON(&jsonGetRP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := jsonGetRP.Email
	code := jsonGetRP.Code
	password := jsonGetRP.Password
	if !email_exist_codegiven(email) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "email",
				},
			},
		})
	} else if !verify_password(password) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "password",
				},
			},
		})
	} else if !code_match(code, email) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "verification_code",
				},
			},
		})
	} else {
		update_password(email, password)
		c.JSON(200, gin.H{})
	}
}

func login(c *gin.Context) {
	var jsonGetL GetLogin
	if err := c.ShouldBindJSON(&jsonGetL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := jsonGetL.Email
	username := jsonGetL.Username
	password := jsonGetL.Password
	if username == "" && email == "" {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "missing_field",
					"field": "email and username",
				},
			},
		})
	} else if username != "" {
		if !username_exist_user(username) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "username",
					},
				},
			})
		} else {
			userId, valid := password_match_username(username, password)
			if !valid {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid",
							"field": "password",
						},
					},
				})
			} else {
				token, _ := utils.GenToken((int64)(userId))
				c.JSON(200, gin.H{
					"access_token": token,
				})
			}
		}
	} else {
		if !email_exist_user(email) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "email",
					},
				},
			})
		} else {
			userId, valid := password_match_email(email, password)
			if !valid {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid",
							"field": "password",
						},
					},
				})
			} else {
				token, _ := utils.GenToken((int64)(userId))
				c.JSON(200, gin.H{
					"access_token": token,
				})
			}
		}
	}
}

func logout(c *gin.Context) {
	token := c.Request.Header["Authorization"][0]
	config.Rdb.SAdd(config.Ctx, "token_blacklist", token[7:])
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
	sqlStr := fmt.Sprintf("select * from codegiven where email='%s';", email)
	//2.执行
	_, ok := config.Query(sqlStr) //从连接池里取一个连接出来去数据库查询记录
	if !ok {
		sqlStr = fmt.Sprintf(`insert into codegiven(email,code,time) values("%s","%s",%d)`, email, code, timeUnix) //sql语句
		ret, err := config.Db.Exec(sqlStr)                                                                         //执行sql语句
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
	} else {
		sqlStr := `update codegiven set code=?,time=? where email=?`
		_, err := config.Db.Exec(sqlStr, code, timeUnix, email)
		if err != nil {
			fmt.Printf("update failed ,err:%v\n", err)
			return
		}
	}
	mail.MailTo(email, code)
}

func code_match(code string, email string) bool {
	//1.查询记录的sql语句
	sqlStr := fmt.Sprintf("select * from codegiven where email='%s';", email)
	//2.执行
	rst, ok := config.Query(sqlStr) //从连接池里取一个连接出来去数据库查询记录
	if ok {
		timeUnix := time.Now().Unix()
		timeThen, _ := strconv.ParseInt(rst[0]["time"], 10, 64)
		return code == rst[0]["code"] && timeUnix-timeThen < 30000 //TODO:change to 300
	} else {
		return false
	}
}

func set_user(email string, password string, username string) int {
	//删除数据库db中的数据
	sqlStr := `delete from codegiven where email=?`
	_, err := config.Db.Exec(sqlStr, email)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return -1
	}

	sqlStr = fmt.Sprintf(`insert into user(username,email,password) values("%s","%s","%s")`, username, email, password) //sql语句
	ret, err := config.Db.Exec(sqlStr)                                                                                  //执行sql语句
	if err != nil {
		fmt.Printf("insert failed,err:%v\n", err)
		return -1
	}
	//如果是插入数据的操作，能够拿到插入数据的id
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get id failed,err:%v\n", err)
		return -1
	}

	return (int)(id)
}

func verify_password(password string) bool {
	expr := `^[0-9A-Za-z!@#$%^&*]{8,16}$`
	reg := regexp.MustCompile(expr)
	m := reg.MatchString(password)
	return m
}

func email_exist_codegiven(email string) bool {
	//1.查询记录的sql语句
	sqlStr := fmt.Sprintf("select * from user where email='%s';", email)
	//2.执行
	_, ok := config.Query(sqlStr) //从连接池里取一个连接出来去数据库查询记录
	return ok
}

func update_email(email string, newEmail string) {
	sqlStr := `update user set email=? where email=?`
	_, err := config.Db.Exec(sqlStr, newEmail, email)
	if err != nil {
		fmt.Printf("update failed ,err:%v\n", err)
		return
	}
}

func update_password(email string, password string) {
	sqlStr := `update user set password=? where email=?`
	_, err := config.Db.Exec(sqlStr, password, email)
	if err != nil {
		fmt.Printf("update failed ,err:%v\n", err)
		return
	}
}

func email_exist_user(email string) bool {
	sqlStr := fmt.Sprintf("select * from user where email='%s';", email)
	_, ok := config.Query(sqlStr)
	return ok
}

func username_exist_user(username string) bool {
	sqlStr := fmt.Sprintf("select * from user where username='%s';", username)
	_, ok := config.Query(sqlStr)
	return ok
}

func password_match_username(username string, password string) (int, bool) {
	sqlStr := fmt.Sprintf("select * from user where username='%s';", username)
	rst, ok := config.Query(sqlStr)
	if ok {
		id, _ := strconv.Atoi(rst[0]["user_id"])
		return id, rst[0]["password"] == password
	} else {
		return 0, false
	}
}

func password_match_email(email string, password string) (int, bool) {
	sqlStr := fmt.Sprintf("select * from user where email='%s';", email)
	rst, ok := config.Query(sqlStr)
	if ok {
		id, _ := strconv.Atoi(rst[0]["user_id"])
		return id, rst[0]["password"] == password
	} else {
		return 0, false
	}
}
