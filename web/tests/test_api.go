// 请在数据库清零的状态下进行测试，或是在完成超级用户后修改
package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"hiper-backend/model"
	"hiper-backend/user"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

var ip string = viper.GetString("OuterIp")
var verificationCode string
var jwtToken string

func init() {
	for i := 0; i < 5; i++ {
		user := model.User{
			Username: fmt.Sprintf("test%d", i),
			Password: user.HashPassword("password"),
			Email:    fmt.Sprintf("test%d@example.com", i),
		}
		err := model.CreateUser(user)
		if err != nil {
			fmt.Printf("Failed to create user: %v", err)
		}
	}
}

func TestWholeUserApi(t *testing.T) {
	t.Run("TestRequestVerificationCode", TestRequestVerificationCode)
	time.Sleep(5 * time.Second)
	t.Run("TestFailRegisterUser", TestFailRegisterUser)
	t.Run("TestRegisterUser", TestRegisterUser)
	t.Run("TestResetPassword", TestResetPassword)
	t.Run("TestResetEmail", TestResetEmail)
	t.Run("TestFailLogin", TestFailLogin)
	t.Run("TestLogin", TestLogin)
	t.Run("TestSearchUsers", TestSearchUsers)
	t.Run("TestGetTheUser", TestGetTheUser)
	t.Run("TestUpdateTheUser", TestUpdateTheUser)
	t.Run("TestGetCurrentUser", TestGetCurrentUser)
}

func TestWholePermissionApi(t *testing.T) {
	t.Run("TestFailGrantCreationPermission", TestFailGrantCreationPermission)
	t.Run("TestGrantCreationPermission", TestGrantCreationPermission)
	t.Run("TestRevokeCreationPermission", TestRevokeCreationPermission)
}

func QuestApi(url string, method string, path string, header map[string]string, query map[string]string, body map[string]string) (map[string]interface{}, error) {
	client := &http.Client{}
	reqURL := url + "/" + path
	reqBody, _ := json.Marshal(body)
	request, err := http.NewRequestWithContext(context.Background(), method, reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range header {
		request.Header.Set(key, value)
	}

	// Set query parameters
	q := request.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	request.URL.RawQuery = q.Encode()

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getLoginToken(email string, password string) {
	login, _ := QuestApi(ip, "POST", "api/v1/user/login", nil, nil, map[string]string{"email": email, "password": password})
	jwtToken = "Bearer " + login["data"].(map[string]interface{})["token"].(string)
}

func TestRequestVerificationCode(t *testing.T) {
	result, err := QuestApi(ip, "POST", "api/v1/user/verification_code", nil, nil, map[string]string{"email": "2975587905@qq.com"})
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	verificationCode, err = user.FindVerificationCode()
	assert.Nil(t, err)
}

func TestFailRegisterUser(t *testing.T) {
	result, err := QuestApi(ip, "POST", "api/v1/users", nil, nil, map[string]string{"username": "test", "password": "lq12345",
		"email": "2975587905@qq.com", "verification_code": verificationCode})
	assert.NotNil(t, err)
	assert.Equal(t, 422, result["code"])
}

func TestRegisterUser(t *testing.T) {
	result, err := QuestApi(ip, "POST", "api/v1/users", nil, nil, map[string]string{"username": "test", "password": "Lq12345",
		"email": "2975587905@qq.com", "verification_code": verificationCode})
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	assert.Equal(t, "test", result["data"].(map[string]interface{})["username"])
}

func TestResetPassword(t *testing.T) {
	result, err := QuestApi(ip, "PUT", "api/v1/user/password", nil,
		nil, map[string]string{"email": "2975587905@qq.com", "verification_code": verificationCode, "new_password": "Lq23456"})
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
}

func TestResetEmail(t *testing.T) {
	result, err := QuestApi(ip, "PUT", "api/v1/user/email", nil,
		nil, map[string]string{"email": "2975587905@qq.com", "verification_code": verificationCode, "new_email": "22@qq.com"})
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
}

func TestLogin(t *testing.T) {
	result, err := QuestApi(ip, "POST", "api/v1/user/login", nil, nil, map[string]string{"email": "22@qq.com", "password": "Lq23456"})
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	jwtToken = "Bearer " + result["data"].(map[string]interface{})["token"].(string)
}

func TestFailLogin(t *testing.T) {
	result, err := QuestApi(ip, "POST", "api/v1/user/login", nil, nil, map[string]string{"email": "22@qq.com"})
	assert.NotNil(t, err)
	assert.Equal(t, 422, result["code"])
}

func TestSearchUsers(t *testing.T) {
	result, err := QuestApi(ip, "GET", "api/v1/users", map[string]string{"Authorization": jwtToken},
		map[string]string{"email": "test1@example.com"}, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	assert.Equal(t, "test1", result["data"].([]interface{})[0].(map[string]interface{})["username"])
}

func TestGetTheUser(t *testing.T) {
	result, err := QuestApi(ip, "GET", "api/v1/users/test2", map[string]string{"Authorization": jwtToken}, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	assert.Equal(t, "test2@example.com", result["data"].(map[string]interface{})["email"])
}

func TestUpdateTheUser(t *testing.T) {
	result, err := QuestApi(ip, "PATCH", "api/v1/user", map[string]string{"Authorization": jwtToken}, nil,
		map[string]string{"username": "testnow", "bio": "test bio"})
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	assert.Equal(t, "testnow", result["data"].(map[string]interface{})["username"])
	assert.Equal(t, "test bio", result["data"].(map[string]interface{})["bio"])
}

func TestGetCurrentUser(t *testing.T) {
	result, err := QuestApi(ip, "GET", "api/v1/user", map[string]string{"Authorization": jwtToken}, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	assert.Equal(t, "testnow", result["data"].(map[string]interface{})["username"])
	assert.Equal(t, "test bio", result["data"].(map[string]interface{})["bio"])
}

func TestFailGrantCreationPermission(t *testing.T) {
	getLoginToken("test3@example.com", "password")
	result, err := QuestApi(ip, "PUT", "api/v1/permissions/create_game_or_contest/2", map[string]string{"Authorization": jwtToken}, nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, 401, result["code"])
}

func TestGrantCreationPermission(t *testing.T) {
	getLoginToken("test0@example.com", "password")
	result, err := QuestApi(ip, "PUT", "api/v1/permissions/create_game_or_contest/2", map[string]string{"Authorization": jwtToken}, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	result, err = QuestApi(ip, "GET", "api/v1/users/test2", map[string]string{"Authorization": jwtToken}, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	assert.Equal(t, true, result["data"].(map[string]interface{})["can_create_game_or_contest"])
}

func TestRevokeCreationPermission(t *testing.T) {
	getLoginToken("test0@example.com", "password")
	result, err := QuestApi(ip, "DELETE", "api/v1/permissions/create_game_or_contest/2", map[string]string{"Authorization": jwtToken}, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	result, err = QuestApi(ip, "GET", "api/v1/users/test2", map[string]string{"Authorization": jwtToken}, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["code"])
	assert.Equal(t, false, result["data"].(map[string]interface{})["can_create_game_or_contest"])
}
