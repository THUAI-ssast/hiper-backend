// 请在数据库清零的状态下进行测试，在完成超级用户后修改
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hiper-backend/api"
	"hiper-backend/config"
	"hiper-backend/model"
	"hiper-backend/user"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var jwtToken string

func makeRequest(url, method, payload string, useAuthorization bool) (map[string]interface{}, error) {
	body := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if useAuthorization {
		req.Header.Add("Authorization", jwtToken)
	}

	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&result)
	result["status"] = res.StatusCode
	if err != nil {
		return nil, err
	}

	return result, nil
}

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
		err := model.CreateUser(&user)
		if err != nil {
			fmt.Printf("Failed to create user: %v", err)
		}
	}
	model.SaveVerificationCode("999999", "2975587905@qq.com", 5)
}

func TestWholeUserApi(t *testing.T) {
	t.Run("TestRequestVerificationCode", TestRequestVerificationCode)
	t.Run("TestRegisterUser", TestRegisterUser)
	t.Run("TestResetPassword", TestResetPassword)
	t.Run("TestLogin", TestLogin)
	t.Run("TestGetTheUser", TestGetTheUser)
	t.Run("TestUpdateTheUser", TestUpdateTheUser)
	t.Run("TestGetCurrentUser", TestGetCurrentUser)
	t.Run("TestResetEmail", TestResetEmail)
}

func TestWholeGameApi(t *testing.T) {
	t.Run("TestCreateGame", TestCreateGame)
	t.Run("TestUpdateMetaData", TestUpdateMetaData)
	t.Run("TestUpdateContestScript", TestUpdateContestScript)
	t.Run("TestUpdateStates", TestUpdateStates)
	t.Run("TestUpdateGameLogic", TestUpdateGameLogic)
	t.Run("TestUpdateGameDetail", TestUpdateGameDetail)
	t.Run("TestGetSettings", TestGetSettings)
	t.Run("TestGetGame", TestGetGame)
	//t.Run("TestAddSdk", TestAddSdk)
	//t.Run("TestUpdateSdk", TestUpdateSdk)
	//t.Run("TestGetSdk", TestGetSdk)
	t.Run("TestCommitAI", TestCommitAI)
	t.Run("TestEditAINote", TestEditAINote)
	t.Run("TestGetAI", TestGetAI)
	//t.Run("TestGetContestant", TestGetContestant)
	//t.Run("TestAssignAI", TestAssignAI)
	//t.Run("TestRevokeAI", TestRevokeAI)
	//t.Run("TestDeleteSdk", TestDeleteSdk)
	t.Run("TestForkGame", TestForkGame)
	t.Run("TestDeleteGame", TestDeleteGame)
	t.Run("TestAddAdmin", TestAddAdmin)
	t.Run("TestRelinquishAdmin", TestRelinquishAdmin)
}

// func TestWholeContestApi(t *testing.T) {
// 	t.Run("TestCreateContest", TestCreateContest)
// 	t.Run("TestRegisterContest", TestRegisterContest)
// 	t.Run("TestExitContest", TestExitContest)
// }

func TestWholePermissionApi(t *testing.T) {
	t.Run("TestGrantCreationPermission", TestGrantCreationPermission)
	t.Run("TestRevokeCreationPermission", TestRevokeCreationPermission)
}

func getLoginToken(email string, password string) {
	url := "http://localhost:8080/api/v1/user/login"
	method := "POST"
	payload := fmt.Sprintf(`{
        "password": "%s",
        "email": "%s"
    }`, password, email)

	result, err := makeRequest(url, method, payload, false)
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

	result, err := makeRequest(url, method, string(jsonStr), false)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	model.SaveVerificationCode("999999", "2975587905@qq.com", 5)
}

func TestRegisterUser(t *testing.T) {
	url := "http://localhost:8080/api/v1/users"
	method := "POST"
	payload := `{
    "password": "Lq3525926",
    "verification_code": "999999",
    "email": "2975587905@qq.com",
    "username": "test"
}`

	result, err := makeRequest(url, method, payload, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, "test", result["username"].(string))
	assert.Equal(t, 200, result["status"].(int))
}

func TestResetPassword(t *testing.T) {
	url := "http://localhost:8080/api/v1/user/reset-password"
	method := "POST"
	payload := `{
    "email": "2975587905@qq.com",
    "verification_code": "999999",
    "new_password":  "Lq234567"
}`

	result, err := makeRequest(url, method, payload, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestLogin(t *testing.T) {
	url := "http://localhost:8080/api/v1/user/login"
	method := "POST"
	payload := `{
    "password": "Lq234567",
    "email": "2975587905@qq.com"
}`

	result, err := makeRequest(url, method, payload, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	jwtToken = "Bearer " + result["access_token"].(string)
}

func TestGetTheUser(t *testing.T) {
	url := "http://localhost:8080/api/v1/users/test2?fields="
	method := "GET"

	result, err := makeRequest(url, method, "", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	assert.Equal(t, "test2@example.com", result["email"].(string))
}

func TestUpdateTheUser(t *testing.T) {
	url := "http://localhost:8080/api/v1/user"
	method := "PATCH"
	payload := `{
    "bio": "test bio"
}`

	result, err := makeRequest(url, method, payload, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	assert.Equal(t, "test bio", result["bio"].(string))
}

func TestGetCurrentUser(t *testing.T) {
	url := "http://localhost:8080/api/v1/user"
	method := "GET"

	result, err := makeRequest(url, method, "", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	assert.Equal(t, "test bio", result["bio"].(string))
}

func TestResetEmail(t *testing.T) {
	url := "http://localhost:8080/api/v1/user/reset-email"
	method := "POST"
	payload := `{
		"email": "2975587905@qq.com",
		"verification_code": "999999",
		"new_email": "abc@mails.tsinghua.edu.cn"
}`
	result, err := makeRequest(url, method, payload, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestGrantCreationPermission(t *testing.T) {
	getLoginToken("admin@mails.tsinghua.edu.cn", "password")
	url := "http://localhost:8080/api/v1/permissions/create_game_or_contest/2"
	method := "PUT"

	result, err := makeRequest(url, method, "", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))

	url = "http://localhost:8080/api/v1/users/test0?fields="
	method = "GET"

	result, err = makeRequest(url, method, "", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	data := result["permissions"].(map[string]interface{})
	assert.Equal(t, true, data["can_create_game_or_contest"].(bool))
}

func TestRevokeCreationPermission(t *testing.T) {
	getLoginToken("admin@mails.tsinghua.edu.cn", "password")
	url := "http://localhost:8080/api/v1/permissions/create_game_or_contest/2"
	method := "DELETE"

	result, err := makeRequest(url, method, "", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))

	url = "http://localhost:8080/api/v1/users/test0?fields="
	method = "GET"

	result, err = makeRequest(url, method, "", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	data := result["permissions"].(map[string]interface{})
	assert.Equal(t, false, data["can_create_game_or_contest"].(bool))
}

func TestCreateGame(t *testing.T) {
	getLoginToken("admin@mails.tsinghua.edu.cn", "password")
	url := "http://localhost:8080/api/v1/games"
	method := "POST"
	payload := `{
	"new_admin_username": "test1"
}`

	result, err := makeRequest(url, method, payload, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestUpdateMetaData(t *testing.T) {
	getLoginToken("test1@example.com", "password")
	url := "http://localhost:8080/api/v1/games/2"
	method := "PATCH"
	payload := `{
	"readme": "test readme"
}`

	result, err := makeRequest(url, method, payload, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestUpdateContestScript(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/2"
	method := "PATCH"
	payload := `{
	"contest_script": "test contest script"
}`

	result, err := makeRequest(url, method, payload, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestUpdateStates(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/2"
	method := "PATCH"
	payload := `{
	"assign_ai_enabled": true
}`

	result, err := makeRequest(url, method, payload, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestUpdateGameLogic(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/2/game_logic"
	method := "PATCH"
	payload := `{
	"build_game_logic_dockerfile": "test dockerfile"
}`

	result, err := makeRequest(url, method, payload, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestUpdateGameDetail(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/2/match_detail"
	method := "PATCH"
	payload := `{
	"template": "test template"
}`

	result, err := makeRequest(url, method, payload, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestGetSettings(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/2/settings"
	method := "GET"

	result, err := makeRequest(url, method, "", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	game_assets := result["game_assets"].(map[string]interface{})
	match_detail := game_assets["match_detail"].(map[string]interface{})
	assert.Equal(t, "test template", match_detail["template"].(string))
	game_logic := game_assets["game_logic"].(map[string]interface{})
	build_game_logic := game_logic["build_game_logic"].(map[string]interface{})
	assert.Equal(t, "test dockerfile", build_game_logic["dockerfile"].(string))
	states := result["states"].(map[string]interface{})
	assert.Equal(t, true, states["assign_ai_enabled"].(bool))
	contest_assets := result["contest_assets"].(map[string]interface{})
	assert.Equal(t, "test contest script", contest_assets["contest_script"].(string))
	metadata := result["metadata"].(map[string]interface{})
	assert.Equal(t, "test readme", metadata["readme"].(string))
}

func TestGetGame(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/2"
	method := "GET"

	result, err := makeRequest(url, method, "", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	assert.Equal(t, "admin", result["my_privilege"].(string))
}

func TestCommitAI(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/1/ais"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open("ai.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	part1, err := writer.CreateFormFile("ai", filepath.Base("ai.txt"))
	if err != nil {
		log.Fatalf("Failed to create form file: %v", err)
	}

	_, err = io.Copy(part1, file)
	if err != nil {
		log.Fatalf("Failed to copy: %v", err)
	}

	_ = writer.WriteField("note", "test note")
	_ = writer.WriteField("sdk_id", "1")

	err = writer.Close()
	if err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("Authorization", jwtToken)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Connection", "keep-alive")

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	assert.Nil(t, err)
	defer res.Body.Close()

	result := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&result)
	result["status"] = res.StatusCode
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestEditAINote(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/1/ais/1/note"
	method := "PUT"
	payload := `{
	"note": "test note"
}`
	result, err := makeRequest(url, method, payload, true)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestGetAI(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/1/ais/1"
	method := "GET"
	result, err := makeRequest(url, method, "", true)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	assert.Equal(t, "test note", result["note"].(string))
}

func TestForkGame(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/1/fork"
	method := "POST"
	payload := `{
	"new_admin_username": "test2"
}`
	result, err := makeRequest(url, method, payload, true)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	game_id := result["id"].(int)

	url = "http://localhost:8080/api/v1/games/" + strconv.Itoa(game_id)
	method = "GET"
	result, err = makeRequest(url, method, "", true)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestDeleteGame(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/1"
	method := "DELETE"
	result, err := makeRequest(url, method, "", true)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
}

func TestAddAdmin(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/2/admins"
	method := "POST"
	payload := `{
	"new_admin_username": "test3"
}`
	result, err := makeRequest(url, method, payload, true)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))

	getLoginToken("test3@example.com", "password")
	url = "http://localhost:8080/api/v1/games/2"
	method = "GET"
	result, err = makeRequest(url, method, "", true)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))
	assert.Equal(t, "admin", result["my_privilege"].(string))
}

func TestRelinquishAdmin(t *testing.T) {
	url := "http://localhost:8080/api/v1/games/2/admin"
	method := "DELETE"
	result, err := makeRequest(url, method, "", true)
	assert.Nil(t, err)
	assert.Equal(t, 200, result["status"].(int))

	url = "http://localhost:8080/api/v1/games/2"
	method = "GET"
	result, err = makeRequest(url, method, "", true)
	assert.Nil(t, err)
	assert.Equal(t, 401, result["status"].(int))
}
