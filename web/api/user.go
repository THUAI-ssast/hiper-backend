package api

import (
	"fmt"
	"hiper-backend/mail"
	"hiper-backend/model"

	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	Avatar_url  string          `json:"avatar_url"`
	Bio         string          `json:"bio"`
	Department  string          `json:"department"`
	Name        string          `json:"name"`
	Permissions map[string]bool `json:"permissions"`
	School      string          `json:"school"`
	Username    string          `json:"username"`
}

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

type GetUsers struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	User_id  int    `json:"user_id"`
}

type GetUserInfo struct {
	Avatar_url string `json:"avatar_url"`
	Bio        string `json:"bio"`
	Department string `json:"department"`
	Name       string `json:"name"`
	School     string `json:"school"`
	Username   string `json:"username"`
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
		userID := set_user(email, password, username)
		if userID == -1 {
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
				"username": username,
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
			userID, valid := password_match_username(username, password)
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
				token, _ := model.GenToken((int64)(userID))
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
			userID, valid := password_match_email(email, password)
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
				token, _ := model.GenToken((int64)(userID))
				c.JSON(200, gin.H{
					"access_token": token,
				})
			}
		}
	}
}

func logout(c *gin.Context) {
	token := c.Request.Header["Authorization"][0]
	model.SetEX(token[7:], 1, 24)
}

func get_user_search(conditions map[string]interface{}, userSearch *[]User) bool {
	rst, valid := model.SelectMySql("user", conditions)
	if !valid {
		return false
	}
	for i := 0; i < len(rst); i++ {
		*userSearch = append(*userSearch, User{
			Avatar_url: rst[i]["avatar_url"],
			Bio:        rst[i]["bio"],
			Department: rst[i]["department"],
			Name:       rst[i]["name"],
			Permissions: map[string]bool{
				"can_create_game_or_contest": rst[i]["authorization"] != "Regular user",
			},
			School:   rst[i]["school"],
			Username: rst[i]["username"],
		})
	}
	return true
}

func search_users(c *gin.Context) {
	var jsonGetSU GetUsers
	var userSearch []User
	if err := c.ShouldBindJSON(&jsonGetSU); err != nil {
		get_user_search(map[string]interface{}{}, &userSearch)
		c.JSON(200, gin.H{
			"answer": userSearch,
		})
		return
	}
	email := jsonGetSU.Email
	username := jsonGetSU.Username
	userID := jsonGetSU.User_id
	if email == "" && username == "" && userID == 0 {
		get_user_search(map[string]interface{}{}, &userSearch)
		c.JSON(200, gin.H{
			"answer": userSearch,
		})
	} else if email != "" {
		get_user_search(map[string]interface{}{"email": email}, &userSearch)
		c.JSON(200, gin.H{
			"answer": userSearch,
		})
	} else if username != "" {
		get_user_search(map[string]interface{}{"username": username}, &userSearch)
		c.JSON(200, gin.H{
			"answer": userSearch,
		})
	} else {
		get_user_search(map[string]interface{}{"user_id": userID}, &userSearch)
		c.JSON(200, gin.H{
			"answer": userSearch,
		})
	}
}

func get_the_user(c *gin.Context, username string) {
	rst, valid := model.SelectMySql("user", map[string]interface{}{"username": username})
	if !valid {
		c.JSON(404, gin.H{})
	} else {
		c.JSON(200, gin.H{
			"avatar_url": rst[0]["avatar_url"],
			"bio":        rst[0]["bio"],
			"department": rst[0]["department"],
			"name":       rst[0]["name"],
			"permissions": map[string]bool{
				"can_create_game_or_contest": rst[0]["authorization"] != "Regular user",
			},
			"school":              rst[0]["school"],
			"username":            rst[0]["username"],
			"email":               rst[0]["email"],
			"contests_registered": rst[0]["contests_id"], //TODO:USE contest id to get contest_registered
		})
	}
}

func get_current_user(c *gin.Context) {
	userIDs, _ := c.Get("userID")
	userID, _ := userIDs.(int)
	rst, valid := model.SelectMySql("user", map[string]interface{}{"user_id": userID})
	if !valid {
		c.JSON(404, gin.H{})
	} else {
		c.JSON(200, gin.H{
			"avatar_url": rst[0]["avatar_url"],
			"username":   rst[0]["username"],
			"bio":        rst[0]["bio"],
			"department": rst[0]["department"],
			"name":       rst[0]["name"],
			"permissions": map[string]bool{
				"can_create_game_or_contest": rst[0]["authorization"] != "Regular user",
			},
			"school":              rst[0]["school"],
			"email":               rst[0]["email"],
			"contests_registered": rst[0]["contests_id"],
		})
	}
}

func update_current_user(c *gin.Context) {
	userIDs, _ := c.Get("userID")
	userID, _ := userIDs.(int)
	_, valid := model.SelectMySql("user", map[string]interface{}{"user_id": userID})
	if !valid {
		c.JSON(404, gin.H{})
	} else {
		var jsonGetCU GetUserInfo
		if err := c.ShouldBindJSON(&jsonGetCU); err != nil {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "json",
					},
				},
			})
			return
		}

		if username_exist_user(jsonGetCU.Username) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "username",
					},
				},
			})
			return
		}

		if !isValidURL(jsonGetCU.Avatar_url) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "avatar_url",
					},
				},
			})
			return
		}

		updates := map[string]interface{}{
			"avatar_url": jsonGetCU.Avatar_url,
			"username":   jsonGetCU.Username,
			"bio":        jsonGetCU.Bio,
			"department": jsonGetCU.Department,
			"name":       jsonGetCU.Name,
			"school":     jsonGetCU.School,
		}

		for key, value := range updates {
			if len(value.(string)) > 100 { // assuming 100 is the maximum length
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid_too_long",
							"field": key,
						},
					},
				})
				return
			}
		}

		if len(updates) > 0 {
			if !model.UpdateMySQL("user", updates, map[string]interface{}{"user_id": userID}) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
				return
			}
		}
		rst, _ := model.SelectMySql("user", map[string]interface{}{"user_id": userID})
		c.JSON(200, gin.H{
			"avatar_url": rst[0]["avatar_url"],
			"username":   rst[0]["username"],
			"bio":        rst[0]["bio"],
			"department": rst[0]["department"],
			"name":       rst[0]["name"],
			"permissions": map[string]bool{
				"can_create_game_or_contest": rst[0]["authorization"] != "Regular user",
			},
			"school":              rst[0]["school"],
			"email":               rst[0]["email"],
			"contests_registered": rst[0]["contests_id"],
		})
	}
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
	_, valid := model.SelectMySql("codegiven", map[string]interface{}{"email": email})
	if !valid {
		ret, valid := model.InsertMySql("codegiven", map[string]interface{}{"email": email, "code": code, "time": timeUnix})
		if !valid {
			return
		}
		//如果是插入数据的操作，能够拿到插入数据的id
		_, err := ret.LastInsertId()
		if err != nil {
			fmt.Printf("get id failed,err:%v\n", err)
			return
		}
	} else {
		if !model.UpdateMySQL("codegiven", map[string]interface{}{"code": code, "time": timeUnix}, map[string]interface{}{"email": email}) {
			return
		}
	}
	mail.MailTo(email, code)
}

func code_match(code string, email string) bool {
	rst, valid := model.SelectMySql("codegiven", map[string]interface{}{"email": email}) //从连接池里取一个连接出来去数据库查询记录
	if valid {
		timeUnix := time.Now().Unix()
		timeThen, _ := strconv.ParseInt(rst[0]["time"], 10, 64)
		return code == rst[0]["code"] && timeUnix-timeThen < 30000 //change to 300 To test timeout
	} else {
		return false
	}
}

func set_user(email string, password string, username string) int {
	//删除数据库db中的数据
	_, valid := model.DeleteMySQL("codegiven", map[string]interface{}{"email": email})
	if !valid {
		return -1
	}

	ret, valid := model.InsertMySql("user", map[string]interface{}{"email": email, "password": password, "username": username})
	if !valid {
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

func isValidURL(urlStr string) bool {
	if urlStr == "" {
		return true
	}
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

func verify_password(password string) bool {
	expr := `^[0-9A-Za-z!@#$%^&*]{8,16}$`
	reg := regexp.MustCompile(expr)
	m := reg.MatchString(password)
	return m
}

func email_exist_codegiven(email string) bool {
	_, valid := model.SelectMySql("user", map[string]interface{}{"email": email})
	return valid
}

func update_email(email string, newEmail string) {
	model.UpdateMySQL("user", map[string]interface{}{"email": newEmail}, map[string]interface{}{"email": email})
}

func update_password(email string, password string) {
	model.UpdateMySQL("user", map[string]interface{}{"password": password}, map[string]interface{}{"email": email})
}

func email_exist_user(email string) bool {
	_, valid := model.SelectMySql("user", map[string]interface{}{"email": email})
	return valid
}

func username_exist_user(username string) bool {
	_, valid := model.SelectMySql("user", map[string]interface{}{"username": username})
	return valid
}

func password_match_username(username string, password string) (int, bool) {
	rst, valid := model.SelectMySql("user", map[string]interface{}{"username": username})
	if valid {
		id, _ := strconv.Atoi(rst[0]["user_id"])
		return id, rst[0]["password"] == password
	} else {
		return 0, false
	}
}

func password_match_email(email string, password string) (int, bool) {
	rst, valid := model.SelectMySql("user", map[string]interface{}{"email": email})
	if valid {
		id, _ := strconv.Atoi(rst[0]["user_id"])
		return id, rst[0]["password"] == password
	} else {
		return 0, false
	}
}
