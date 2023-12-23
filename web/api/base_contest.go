package api

import (
	"hiper-backend/game"
	"hiper-backend/model"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func deleteGame(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	err := model.DeleteGameByID(gameID)
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
	game, err := model.GetGameByID(gameID)
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
	game, err := model.GetGameByID(gameID)
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
			err = model.DeleteGameByID(gameID)
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
	err := model.UpdateGameByID(gameID, map[string]interface{}{"script": input.Script})
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
	err := model.UpdateGameByID(gameID, map[string]interface{}{
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
		Name              string                `json:"name"`
		Description       string                `json:"description"`
		Sdk               *multipart.FileHeader `json:"sdk"`
		BuildAiDockerfile string                `json:"build_ai_dockerfile"`
		RunAiDockerfile   string                `json:"run_ai_dockerfile"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
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
	err := model.UpdateGameByID(gameID, updates)
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

func getGames(c *gin.Context) {
	games, err := model.GetGames()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	var gamesList []gin.H
	for _, game := range games {
		userID := c.MustGet("userID").(int)
		pri, err := game.GetPrivilege(uint(userID))
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		gameData := gin.H{
			"id":           game.ID,
			"game_id":      game.BaseContest.GameID,
			"metadata":     game.Metadata,
			"states":       game.BaseContest.States,
			"my_privilege": pri,
		}
		gamesList = append(gamesList, gameData)
	}

	c.JSON(200, gamesList)
}

func getTheGame(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	game, err := model.GetGameByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Game not found"})
		return
	}

	userID := c.MustGet("userID").(int)
	pri, err := game.GetPrivilege(uint(userID))
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"id":           game.ID,
		"game_id":      game.BaseContest.GameID,
		"metadata":     game.Metadata,
		"states":       game.BaseContest.States,
		"my_privilege": pri,
		// TODO: my
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

	// file, err := c.FormFile("ai")
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

	// TODO: 上传文件、更新数据库
	// update???

	c.JSON(http.StatusOK, gin.H{"id": gameID})
}

func getTheAI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	_, err = model.GetGameByID(uint(id))
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
	// 	id, err := strconv.Atoi(c.Param("id"))
	// 	if err != nil {
	// 		c.JSON(400, gin.H{})
	// 		return
	// 	}
	// 	_, err = model.GetGameByID(uint(id))
	// 	if err != nil {
	// 		c.JSON(404, gin.H{"error": "Game not found"})
	// 		return
	// 	}

	// 	ai_id, err := strconv.Atoi(c.Param("ai_id"))
	// 	if err != nil {
	// 		c.JSON(400, gin.H{})
	// 		return
	// 	}
	// 	ai, err := model.GetAiByID(uint(ai_id), true)
	// 	if err != nil {
	// 		c.JSON(404, gin.H{"error": "AI not found"})
	// 		return
	// 	}
	// 	//file, err := ai.GetFile()
	// 	if err != nil {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
	// 		return
	// 	}

	// 	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	// 	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "desired_filename.ext"))
	// 	c.Writer.WriteHeader(http.StatusOK)
	// 	c.Writer.Write(file)
}

func editAiNote(c *gin.Context) {
	//aiID, err := strconv.Atoi(c.Param("ai_id"))
	_, err := strconv.Atoi(c.Param("ai_id"))
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

	// TODO: somgthing wrong need to revise
	// err = update(map[string]interface{}{"note": requestBody.Note}, map[string]interface{}{"note": requestBody.Note})
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
	// 	id, err := strconv.Atoi(c.Param("id"))
	// 	if err != nil {
	// 		c.JSON(400, gin.H{})
	// 		return
	// 	}

	// 	game, err := model.GetGameByID(uint(id))
	// 	if err != nil {
	// 		c.JSON(404, gin.H{})
	// 		return
	// 	}

	// 	contestants, err := game.GetContestants()
	// 	if err != nil {
	// 		c.JSON(404, gin.H{})
	// 		return
	// 	}

	// 	var contestantList []gin.H
	// 	for _, contestant := range contestants {
	// 		userid := contestant.UserId
	// 		user, err := model.GetUserByID(uint(userid))
	// 		if err != nil {
	// 			c.JSON(400, gin.H{})
	// 			return
	// 		}

	// 		aiid := contestant.AssignedAiId
	// 		ai, err := game.GetAiByID(uint(aiid), true)
	// 		if err != nil {
	// 			c.JSON(400, gin.H{})
	// 			return
	// 		}

	// 		contestantData := gin.H{
	// 			"performance": contestant.Performance,
	// 			"permissions": contestant.Permissions,
	// 			"points":      contestant.Points,
	// 			"user":        user,
	// 			"assigned_ai": ai,
	// 		}
	// 		contestantList = append(contestantList, contestantData)
	// 	}

	// 	c.JSON(200, contestantList)
}
