// 请在数据库清零的状态下进行测试，在完成超级用户后修改
package tests

import (
	"encoding/json"
	"fmt"
	"hiper-backend/api"
	"hiper-backend/config"
	"hiper-backend/model"
	"hiper-backend/user"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var jwtToken string

func init() {
	config.InitConfig()
	model.InitDb()
	model.AutoMigrateDb()
	model.InitRedis()
	config.InitConfigAfter()
	go api.ApiListenHttp()
	time.Sleep(10 * time.Second)
	for i := 0; i < 5; i++ {
		user := model.User{
			Username: fmt.Sprintf("test%d", i),
			Password: user.HashPassword("password"),
			Email:    fmt.Sprintf("test%d@example.com", i),
		}
		err := user.Create()
		if err != nil {
			fmt.Printf("Failed to create user: %v", err)
		}
	}
	model.SaveVerificationCode("999999", "2975587905@qq.com", 5)
}

func TestWholeUserApi(t *testing.T) {
	t.Run("TestRequestVerificationCode", TestRequestVerificationCode)
	time.Sleep(5 * time.Second)
	t.Run("TestRegisterUser", TestRegisterUser)
	t.Run("TestResetPassword", TestResetPassword)
	t.Run("TestLogin", TestLogin)
	t.Run("TestGetTheUser", TestGetTheUser)
	t.Run("TestUpdateTheUser", TestUpdateTheUser)
	t.Run("TestGetCurrentUser", TestGetCurrentUser)
}

func TestWholePermissionApi(t *testing.T) {
	t.Run("TestGrantCreationPermission", TestGrantCreationPermission)
}

func getLoginToken(email string, password string) {
	url := "http://localhost:8080/api/v1/user/login"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
		"password": "%s",
		"email": "%s"
	}`, password, email))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	jwtToken = "Bearer " + result["access_token"].(string)
}

func TestRequestVerificationCode(t *testing.T) {
	url := "http://localhost:8080/api/v1/user/request-verification-code"
	method := "POST"

	data := map[string]string{
		"email": "2975587905@qq.com",
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	payload := strings.NewReader(string(jsonStr))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Nil(t, err)
	model.SaveVerificationCode("999999", "2975587905@qq.com", 5)
}

func TestRegisterUser(t *testing.T) {
	url := "http://localhost:8080/api/v1/users"
	method := "POST"

	payload := strings.NewReader(`{
    "password": "Lq3525926",
    "verification_code": "999999",
    "email": "2975587905@qq.com",
    "username": "test"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	result := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&result)

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "test", result["username"].(string))
}

func TestResetPassword(t *testing.T) {
	url := "http://localhost:8080/api/v1/user/reset-password"
	method := "POST"

	payload := strings.NewReader(`{
    "email": "2975587905@qq.com",
    "verification_code": "999999",
    "new_password":  "Lq234567"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestLogin(t *testing.T) {
	url := "http://localhost:8080/api/v1/user/login"
	method := "POST"

	payload := strings.NewReader(`{
    "password": "Lq234567",
    "email": "2975587905@qq.com"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	jwtToken = "Bearer " + result["access_token"].(string)
}

func TestGetTheUser(t *testing.T) {
	url := "http://localhost:8080/api/v1/users/test2?fields="
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", jwtToken)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "test2@example.com", result["email"].(string))
}

func TestUpdateTheUser(t *testing.T) {
	url := "http://localhost:8080/api/v1/user"
	method := "PATCH"

	payload := strings.NewReader(`{
    "bio": "test bio"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", jwtToken)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "test bio", result["bio"].(string))
}

func TestGetCurrentUser(t *testing.T) {
	url := "http://localhost:8080/api/v1/user"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", jwtToken)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "test bio", result["bio"].(string))
}

func TestGrantCreationPermission(t *testing.T) {
	//TODO:Change To super admin
	getLoginToken("test0@example.com", "password")
	url := "http://localhost:8080/api/v1/permissions/create_game_or_contest/2"
	method := "PUT"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", jwtToken)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	url = "http://localhost:8080/api/v1/users/test2?fields="
	method = "GET"

	client = &http.Client{}
	req, err = http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", jwtToken)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	res, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, true, result["can_create_game_or_contest"].(bool))
}
