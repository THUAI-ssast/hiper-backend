package api

import (
	"fmt"
	"hiper-backend/game"
	"hiper-backend/model"
	"hiper-backend/mq"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func deleteGame(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	err := model.DeleteBaseContestByID(gameID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()
}

func addAdmin(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	var input struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	user, err := model.GetUserByUsername(input.Username)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  MissingField,
			Field: "user",
		}})
		c.Abort()
		return
	}
	game, err := model.GetBaseContestByID(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  MissingField,
			Field: "game",
		}})
		c.Abort()
		return
	}
	err = game.AddAdmin(user.ID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot add admin",
		}})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()
}

func relinquishAdmin(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	var input struct {
		Force bool `json:"force"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	inuserID := c.MustGet("userID").(int)
	userID := uint(inuserID)
	game, err := model.GetBaseContestByID(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  MissingField,
			Field: "game",
		}})
		c.Abort()
		return
	}
	if admins, _ := game.GetAdmins(); len(admins) == 1 {
		if input.Force {
			err = model.DeleteBaseContestByID(gameID)
			if err != nil {
				c.JSON(422, gin.H{"error": ErrorFor422{
					Code:  Invalid,
					Field: "cannot delete game",
				}})
				c.Abort()
				return
			}
		} else {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "cannot relinquish the only admin",
			}})
			c.Abort()
			return
		}
	}
	err = game.RemoveAdmin(userID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot remove admin",
		}})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()
}

func updateGameScript(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	var input struct {
		Script string `json:"contest_script"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	err := model.UpdateBaseContestByID(gameID, map[string]interface{}{"script": input.Script})
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update game script",
		}})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()
}

func updateGameMetadata(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	var input struct {
		CoverURL string `json:"cover_url"`
		Readme   string `json:"readme"`
		Title    string `json:"title"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	err := model.UpdateBaseContestByID(gameID, map[string]interface{}{
		"cover_url": input.CoverURL,
		"readme":    input.Readme,
		"title":     input.Title,
	})
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update game metadata",
		}})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()
}

func addSdk(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	var input struct {
		Name              string                `form:"name"`
		Description       string                `form:"description"`
		Sdk               *multipart.FileHeader `form:"sdk"`
		BuildAiDockerfile string                `form:"build_ai_dockerfile"`
		RunAiDockerfile   string                `form:"run_ai_dockerfile"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	sdk := model.Sdk{
		Name:          input.Name,
		Readme:        input.Description,
		BaseContestID: gameID,
	}
	sdk.BuildAi.Dockerfile = input.BuildAiDockerfile
	sdk.RunAi.Dockerfile = input.RunAiDockerfile
	err := sdk.Create()
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot create sdk",
		}})
		c.Abort()
		return
	}
	//saveSdkFile(sdk.ID,input.Sdk)
	//TODO:往sdk中添加，存储文件至/var/hiper/sdks/sdks:id.xxx
	mq.SendBuildSdkMsg(model.Ctx, sdk.ID)
	c.JSON(200, gin.H{})
	c.Abort()
}

func getSdk(c *gin.Context) {
	sdk, err := game.GetSdksFromKnownGame(c)
	if err != nil {
		return
	}
	c.JSON(200, gin.H{
		"id":     sdk.ID,
		"name":   sdk.Name,
		"readme": sdk.Readme,
		"build_ai": map[string]interface{}{
			"dockerfile": sdk.BuildAi.Dockerfile,
			"status": map[string]interface{}{
				"state": sdk.BuildAi.Status.State,
				"msg":   sdk.BuildAi.Status.Msg,
			},
		},
		"run_ai": map[string]interface{}{
			"dockerfile": sdk.RunAi.Dockerfile,
			"status": map[string]interface{}{
				"state": sdk.RunAi.Status.State,
				"msg":   sdk.RunAi.Status.Msg,
			},
		},
	})
	c.Abort()
}

func deleteSdk(c *gin.Context) {
	sdk, err := game.GetSdksFromKnownGame(c)
	if err != nil {
		return
	}
	err = model.DeleteSdkByID(sdk.ID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot delete sdk",
		}})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()
}

func updateSdk(c *gin.Context) {
	sdk, err := game.GetSdksFromKnownGame(c)
	if err != nil {
		return
	}
	var input struct {
		Name              string                `json:"name"`
		Readme            string                `json:"readme"`
		Sdk               *multipart.FileHeader `json:"sdk"`
		BuildAiDockerfile string                `json:"build_ai_dockerfile"`
		RunAiDockerfile   string                `json:"run_ai_dockerfile"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	err = sdk.Update(map[string]interface{}{
		"name":                input.Name,
		"readme":              input.Readme,
		"build_ai_dockerfile": input.BuildAiDockerfile,
		"run_ai_dockerfile":   input.RunAiDockerfile,
	})
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update sdk",
		}})
		c.Abort()
		return
	}
	//saveSdkFile(sdk.ID,input.Sdk)
	//TODO:往sdk中添加，存储文件至/var/hiper/sdks/sdks:id.xxx
	c.JSON(200, gin.H{})
	c.Abort()
}

func updateGameStates(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	var input struct {
		AssignAiEnabled                 *bool `json:"assign_ai_enabled"`
		CommitAiEnabled                 *bool `json:"commit_ai_enabled"`
		ContestScriptEnvironmentEnabled *bool `json:"contest_script_environment_enabled"`
		PrivateMatchEnabled             *bool `json:"private_match_enabled"`
		PublicMatchEnabled              *bool `json:"public_match_enabled"`
		TestMatchEnabled                *bool `json:"test_match_enabled"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	updates := make(map[string]interface{})
	if input.AssignAiEnabled != nil {
		updates["assign_ai_enabled"] = input.AssignAiEnabled
	}
	if input.CommitAiEnabled != nil {
		updates["commit_ai_enabled"] = input.CommitAiEnabled
	}
	if input.ContestScriptEnvironmentEnabled != nil {
		updates["contest_script_environment_enabled"] = input.ContestScriptEnvironmentEnabled
	}
	if input.PrivateMatchEnabled != nil {
		updates["private_match_enabled"] = input.PrivateMatchEnabled
	}
	if input.PublicMatchEnabled != nil {
		updates["public_match_enabled"] = input.PublicMatchEnabled
	}
	if input.TestMatchEnabled != nil {
		updates["test_match_enabled"] = input.TestMatchEnabled
	}
	err := model.UpdateBaseContestByID(gameID, updates)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update game states",
		}})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()
}

func getTheGame(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	baseContest, err := model.GetBaseContestByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Game not found"})
		return
	}

	// baseContest.GetPriviliege(uint(userID))
	admins, err := baseContest.GetAdmins()
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	pri := "registered"
	userID := c.MustGet("userID").(int)
	for _, admin := range admins {
		if admin.ID == uint(userID) {
			pri = "admin"
			break
		}
	}

	preloads := []model.PreloadQuery{
		{
			Table:   "User",
			Columns: []string{},
		},
		{
			Table:   "Ai",
			Columns: []string{},
		},
	}
	contestants, err := baseContest.GetContestants(preloads)
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	var my model.Contestant
	found := false
	for _, contestant := range contestants {
		if contestant.UserID == uint(userID) {
			my = contestant
			found = true
			break
		}
	}
	if !found {
		my = model.Contestant{}
	}

	var input struct {
		CoverURL string `json:"cover_url"`
		Readme   string `json:"readme"`
		Title    string `json:"title"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"base_contest": map[string]interface{}{
			"id":      baseContest.ID,
			"game_id": baseContest.GameID,
			"states":  baseContest.States,
			"my":      my,
		},
		"id":           baseContest.ID,
		"my_privilege": pri,
		"metadata": map[string]interface{}{
			"cover_url": input.CoverURL,
			"readme":    input.Readme,
			"title":     input.Title,
		},
	})
}

func getAis(c *gin.Context) {
	username := c.Query("username")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	queryParams := model.QueryParams{
		Filter: map[string]interface{}{},
		Limit:  limit,
		Offset: offset,
	}
	if username != "" {
		queryParams.Filter["username"] = username
	}

	ais, _, err := model.GetAis(queryParams, true)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	var aiList []gin.H
	for _, ai := range ais {
		aiData := gin.H{
			"id":     ai.ID,
			"sdk":    ai.Sdk,
			"note":   ai.Note,
			"status": ai.Status,
			"user":   ai.User,
			"time":   ai.CreatedAt, // TODO: 可能代表创建时间
		}
		aiList = append(aiList, aiData)
	}

	response := gin.H{
		"count": len(ais),
		"data":  aiList,
	}
	c.JSON(200, response)
}

func commitAi(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []map[string]interface{}{
				{
					"code":   "invalid",
					"field":  "id",
					"detail": "Game ID must be an integer",
				},
			},
		})
		return
	}

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB is the default used by net/http
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []map[string]interface{}{
				{
					"code":   "invalid",
					"detail": "Could not parse multipart form",
				},
			},
		})
		return
	}

	file, err := c.FormFile("ai")
	_, err = c.FormFile("ai")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []map[string]interface{}{
				{
					"code":   "missing_field",
					"field":  "ai",
					"detail": "AI file is required",
				},
			},
		})
		return
	}

	note := c.PostForm("note")
	sdkID := c.PostForm("sdk_id")

	var missingFields []map[string]interface{}
	if note == "" {
		missingFields = append(missingFields, map[string]interface{}{
			"code":   "missing_field",
			"field":  "note",
			"detail": "Note is required",
		})
	}
	if sdkID == "" {
		missingFields = append(missingFields, map[string]interface{}{
			"code":   "missing_field",
			"field":  "sdk_id",
			"detail": "SDK ID is required",
		})
	}
	if len(missingFields) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": missingFields})
		return
	}

	var ai *model.Ai
	ai.Create()
	ai.BaseContestID = uint(gameID)

	model.UpdateAiByID(ai.ID, map[string]interface{}{
		"note":   note,
		"sdk_id": sdkID,
	})

	// TODO:添加到/var/hiper/ais/ais:id
	// Open the uploaded file
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open uploaded file"})
		return
	}
	defer openedFile.Close()

	// Read a chunk to get the file type
	buffer := make([]byte, 512)
	if _, err = openedFile.Read(buffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read uploaded file"})
		return
	}

	// Detect content type
	contentType := http.DetectContentType(buffer)

	// Reset the read pointer of the file
	if _, err = openedFile.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reset file read pointer"})
		return
	}

	// Map MIME type to file extension
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || len(exts) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to determine file extension"})
		return
	}

	// Use the first extension returned by the mime package
	extension := exts[0]

	// Construct file path
	aiFilePath := fmt.Sprintf("/var/hiper/ais/ais:%d%s", ai.ID, extension)

	// Save the file
	if err := c.SaveUploadedFile(file, aiFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to save AI file",
			"detail": err.Error(),
		})
		return
	}

	//以下的两个其实均为AIID，需要修改
	mq.SendBuildAIMsg(model.Ctx, uint(ai.ID))
	c.JSON(http.StatusOK, gin.H{"id": ai.ID})
}

func getTheAI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	_, err = model.GetBaseContestByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Game not found"})
		return
	}

	ai_id, err := strconv.Atoi(c.Param("ai_id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	ai, err := model.GetAiByID(uint(ai_id), true)
	if err != nil {
		c.JSON(404, gin.H{"error": "AI not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":     ai.ID,
		"sdk":    ai.Sdk,
		"note":   ai.Note,
		"user":   ai.User,
		"status": ai.Status,
	})
}

func downloadTheAI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	_, err = model.GetBaseContestByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Game not found"})
		return
	}

	ai_id, err := strconv.Atoi(c.Param("ai_id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	// Construct the base path for the AI file
	aiFilePath := fmt.Sprintf("/var/hiper/ais/ais:%d", ai_id)
	fileDir := "/var/hiper/ais/"
	var fileName string

	// Search for the file with the correct ai_id and extension
	err = filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(info.Name(), aiFilePath) {
			fileName = path
			return nil
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching for file"})
		return
	}

	if fileName == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "AI file not found"})
		return
	}

	// Read the file content
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read AI file"})
		return
	}

	// Send the file as a download
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(fileName)))
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(fileContent)
}

func editAiNote(c *gin.Context) {
	aiID, err := strconv.Atoi(c.Param("ai_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid AI ID",
		})
		return
	}

	// 解析请求体中的新附注
	var requestBody struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err = model.UpdateAiByID(uint(aiID), map[string]interface{}{"note": requestBody.Note})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "AI not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "AI note updated successfully",
	})
}

func getContestants(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	baseContest, err := model.GetBaseContestByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	preloads := []model.PreloadQuery{
		{
			Table:   "User",
			Columns: []string{},
		},
		{
			Table:   "Ai",
			Columns: []string{},
		},
	}
	contestants, err := baseContest.GetContestants(preloads)
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	var contestantList []gin.H
	for _, contestant := range contestants {
		userid := contestant.UserID
		user, err := model.GetUserByID(uint(userid))
		if err != nil {
			c.JSON(400, gin.H{})
			return
		}

		aiid := contestant.AssignedAiID
		ai, err := model.GetAiByID(uint(aiid), true)
		if err != nil {
			c.JSON(400, gin.H{})
			return
		}

		contestantData := gin.H{
			"assigned_ai": ai,
			"performance": contestant.Performance,
			"permissions": contestant.Permissions,
			"points":      contestant.Points,
			"user":        user,
		}
		contestantList = append(contestantList, contestantData)
	}

	c.JSON(200, contestantList)
}

func assignAi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	baseContest, err := model.GetBaseContestByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	preloads := []model.PreloadQuery{
		{
			Table:   "User",
			Columns: []string{},
		},
		{
			Table:   "Ai",
			Columns: []string{},
		},
		{
			Table:   "Contestant",
			Columns: []string{},
		},
	}
	userID := c.MustGet("userID").(int)
	contestant, err := baseContest.GetContestantByUserID(uint(userID), preloads)
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	var input struct {
		AIID int `json:"ai_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	ai_id := input.AIID
	ai, err := model.GetAiByID(uint(ai_id), true)
	if err != nil {
		c.JSON(404, gin.H{"error": "AI not found"})
		return
	}

	contestant.AssignedAi = ai
	contestant.AssignedAiID = uint(ai_id)

	c.JSON(200, gin.H{})
}

func getCurrentContestant(c *gin.Context) {
	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	c.JSON(400, gin.H{})
	// 	return
	// }

	// baseContest, err := model.GetBaseContestByID(uint(id))
	// if err != nil {
	// 	c.JSON(404, gin.H{})
	// 	return
	// }

	contestantIDs, _ := c.Get("contestantID") //
	contestantID, _ := contestantIDs.(int)
	preloads := []model.PreloadQuery{
		{
			Table:   "User",
			Columns: []string{},
		},
		{
			Table:   "Ai",
			Columns: []string{},
		},
		{
			Table:   "Contestant",
			Columns: []string{},
		},
	}
	contestant, err := model.GetContestantByID((uint)(contestantID), preloads)

	aiid := contestant.AssignedAiID
	ai, err := model.GetAiByID(uint(aiid), true)
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	userid := contestant.UserID
	user, err := model.GetUserByID(uint(userid))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"assigned_ai": ai,
		"performance": contestant.Performance,
		"permissions": contestant.Permissions,
		"points":      contestant.Points,
		"user":        user,
	})
}

func revokeAssignedAi(c *gin.Context) {
	contestantIDs, _ := c.Get("contestantID") //
	contestantID, _ := contestantIDs.(int)
	preloads := []model.PreloadQuery{
		{
			Table:   "User",
			Columns: []string{},
		},
		{
			Table:   "Ai",
			Columns: []string{},
		},
		{
			Table:   "Contestant",
			Columns: []string{},
		},
	}
	contestant, err := model.GetContestantByID((uint)(contestantID), preloads)
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	contestant.AssignedAi = model.Ai{}
	contestant.AssignedAiID = 0

	c.JSON(200, gin.H{})
}

func getMatches(c *gin.Context) {
	username := c.Query("username")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	queryParams := model.QueryParams{
		Filter: map[string]interface{}{},
		Offset: offset,
		Limit:  limit,
	}
	if username != "" {
		queryParams.Filter["username"] = username
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	baseContest, err := model.GetBaseContestByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	matches, count, err := baseContest.GetMatches(queryParams, true)
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	var matchList []gin.H
	for _, match := range matches {
		matchData := gin.H{
			"id":      match.ID,
			"tag":     match.Tag,
			"players": match.Players,
			"state":   match.State,
			"time":    match.CreatedAt, // 可能代表创建时间
		}
		matchList = append(matchList, matchData)
	}

	response := gin.H{
		"count": count,
		"data":  matchList,
	}
	c.JSON(200, response)
}

func getMatch(c *gin.Context) {
	username := c.Query("username")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	queryParams := model.QueryParams{
		Filter: map[string]interface{}{},
		Offset: offset,
		Limit:  limit,
	}
	if username != "" {
		queryParams.Filter["username"] = username
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	baseContest, err := model.GetBaseContestByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	matches, _, err := baseContest.GetMatches(queryParams, true)
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	match_id, err := strconv.Atoi(c.Param("match_id"))
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	var aimMatch model.Match
	for _, match := range matches {
		if match.ID == uint(match_id) {
			aimMatch = match
			break
		}
	}

	c.JSON(200, gin.H{
		"id":      aimMatch.ID,
		"tag":     aimMatch.Tag,
		"state":   aimMatch.State,
		"time":    aimMatch.CreatedAt,
		"players": aimMatch.Players,
		"replay":  "", // TODO
	})
}

func getSdks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	baseContest, err := model.GetBaseContestByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Game not found"})
		return
	}

	sdks, err := baseContest.GetSdks()
	if err != nil {
		c.JSON(404, gin.H{})
		return
	}

	var sdkList []gin.H
	for _, sdk := range sdks {
		sdkData := gin.H{
			"id":     sdk.ID,
			"name":   sdk.Name,
			"readme": sdk.Readme,
			"build_ai": map[string]string{
				"status": string(sdk.BuildAi.Status.State),
			},
			"run_ai": map[string]string{
				"status": string(sdk.RunAi.Status.State),
			},
		}
		sdkList = append(sdkList, sdkData)
	}

	c.JSON(200, sdkList)
}
