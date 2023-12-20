package api

import (
	"hiper-backend/game"
	"hiper-backend/model"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createGame(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userr, err := model.GetUserById(userID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find user",
		}})
		return
	}
	if !userr.Permissions.CanCreateGameOrContest {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot create game",
		}})
		return
	}
	var input struct {
		NewAdminUsername string `json:"new_admin_username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tempGame := model.Game{}
	if input.NewAdminUsername != "" {
		newAdmin, err := model.GetUserByUsername(input.NewAdminUsername)
		if err != nil {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "cannot find user",
			}})
			return
		}
		err = model.CreateGame(&tempGame, []uint{newAdmin.ID})
		if err != nil {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "cannot create game",
			}})
			return
		}
	} else {
		err = model.CreateGame(&tempGame, []uint{userID})
	}
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot create game",
		}})
		return
	}
	//TODO:TRUE?
	c.JSON(200, gin.H{"id": tempGame.ID})
}

func forkGame(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userr, err := model.GetUserById(userID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find user",
		}})
		return
	}
	if !userr.Permissions.CanCreateGameOrContest {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot create game",
		}})
		return
	}
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}
	var input struct {
		NewAdminUsername string `json:"new_admin_username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tempGame, err := model.GetGameById(uint(gameID))
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find game",
		}})
		return
	}
	if input.NewAdminUsername != "" {
		newAdmin, err := model.GetUserByUsername(input.NewAdminUsername)
		if err != nil {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "cannot find user",
			}})
			return
		}
		err = model.CreateGame(&tempGame, []uint{newAdmin.ID})
		if err != nil {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "cannot create game",
			}})
			return
		}
	} else {
		err = model.CreateGame(&tempGame, []uint{userID})
	}
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot create game",
		}})
		return
	}
	//TODO:TRUE?
	c.JSON(200, gin.H{"id": tempGame.ID})
}

func deleteGame(c *gin.Context) {
	err := model.DeleteGameById(c.MustGet("gameID").(uint))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{})
}

func addAdmin(c *gin.Context) {
	gameID := c.MustGet("gameID").(uint)
	var input struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user, err := model.GetUserByUsername(input.Username)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  MissingField,
			Field: "user",
		}})
	}
	game, err := model.GetGameById(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  MissingField,
			Field: "game",
		}})
	}
	err = game.AddAdmin(user.ID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot add admin",
		}})
	}
	c.JSON(200, gin.H{})
}

func relinquishAdmin(c *gin.Context) {
	gameID := c.MustGet("gameID").(uint)
	var input struct {
		Force bool `json:"force"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	userID := c.MustGet("userID").(uint)
	game, err := model.GetGameById(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  MissingField,
			Field: "game",
		}})
	}
	if admins, _ := game.GetAdmins(); len(admins) == 1 {
		if input.Force {
			err = model.DeleteGameById(gameID)
			if err != nil {
				c.JSON(422, gin.H{"error": ErrorFor422{
					Code:  Invalid,
					Field: "cannot delete game",
				}})
				return
			}
		} else {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "cannot relinquish the only admin",
			}})
		}
	}
	err = game.RemoveAdmin(userID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot remove admin",
		}})
	}
	c.JSON(200, gin.H{})
}

func updateGameScript(c *gin.Context) {
	gameID := c.MustGet("gameID").(uint)
	var input struct {
		Script string `json:"contest_script"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := model.UpdateGameById(gameID, map[string]interface{}{"script": input.Script})
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update game script",
		}})
	}
	c.JSON(200, gin.H{})
}

func updateGameMetadata(c *gin.Context) {
	gameID := c.MustGet("gameID").(uint)
	var input struct {
		CoverURL string `json:"cover_url"`
		Readme   string `json:"readme"`
		Title    string `json:"title"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := model.UpdateGameById(gameID, map[string]interface{}{
		"cover_url": input.CoverURL,
		"readme":    input.Readme,
		"title":     input.Title,
	})
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update game metadata",
		}})
	}
	c.JSON(200, gin.H{})
}

func addSdk(c *gin.Context) {
	// gameID := c.MustGet("gameID").(uint)
	var input struct {
		Name              string                `json:"name"`
		Description       string                `json:"description"`
		Sdk               *multipart.FileHeader `json:"sdk"`
		BuildAiDockerfile string                `json:"build_ai_dockerfile"`
		RunAiDockerfile   string                `json:"run_ai_dockerfile"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//TODO:往sdk中添加，存储文件至/var/hiper/sdks/sdks:id.xxx
}

func getSdk(c *gin.Context) {
	//TODO:检测sdk是对应game的，然后获取
}

func deleteSdk(c *gin.Context) {
	//TODO:检测sdk是对应game的，然后删除
}

func updateSdk(c *gin.Context) {
	//TODO:检测sdk是对应game的，然后更新
}

func updateGameStates(c *gin.Context) {
	gameID := c.MustGet("gameID").(uint)
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
	err := model.UpdateGameById(gameID, updates)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update game states",
		}})
		return
	}
	c.JSON(200, gin.H{})
}

func getGameSettings(c *gin.Context) {
	game.RetGameSettings(c)
}

func updateGameLogic(c *gin.Context) {
	gameID := c.MustGet("gameID").(uint)
	var input struct {
		BuildGameLogicDockerfile string `json:"build_game_logic_dockerfile"`
		RunGameLogicDockerfile   string `json:"run_game_logic_dockerfile"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := model.UpdateGameById(gameID, map[string]interface{}{
		"build_game_logic_dockerfile": input.BuildGameLogicDockerfile,
		"run_game_logic_dockerfile":   input.RunGameLogicDockerfile,
	})
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update game logic",
		}})
		return
	}
	game.RetGameSettings(c)
}

func updateMatchDetail(c *gin.Context) {
	gameID := c.MustGet("gameID").(uint)
	var input struct {
		Template string `json:"template"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := model.UpdateGameById(gameID, map[string]interface{}{
		"template": input.Template,
	})
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update match detail",
		}})
		return
	}
	c.JSON(200, gin.H{})
}
