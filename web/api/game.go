package api

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/THUAI-ssast/hiper-backend/web/basecontest"
	"github.com/THUAI-ssast/hiper-backend/web/game"
	"github.com/THUAI-ssast/hiper-backend/web/model"
)

func createGame(c *gin.Context) {
	inuserID := c.MustGet("userID").(int)
	userID := uint(inuserID)
	usr, err := model.GetUserByID(userID)
	if err != nil {
		c.JSON(401, gin.H{})
		c.Abort()
		return
	}
	if !usr.Permissions.CanCreateGameOrContest {
		c.JSON(403, gin.H{})
		c.Abort()
		return
	}

	var input struct {
		NewAdminUsername string `json:"new_admin_username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	tempGame := model.Game{}
	if input.NewAdminUsername == "" {
		err = tempGame.Create([]uint{userID})
	} else {
		newAdmin, err := model.GetUserByUsername(input.NewAdminUsername)
		if err != nil {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "cannot find user",
			}})
			c.Abort()
			return
		}
		err = tempGame.Create([]uint{newAdmin.ID})
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create game"})
			c.Abort()
			return
		}
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create game"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"id": tempGame.ID})
	c.Abort()
}

func forkGame(c *gin.Context) {
	inuserID := c.MustGet("userID").(int)
	userID := uint(inuserID)
	userr, err := model.GetUserByID(userID)
	if err != nil {
		c.JSON(401, gin.H{})
		c.Abort()
		return
	}
	if !userr.Permissions.CanCreateGameOrContest {
		c.JSON(403, gin.H{})
		c.Abort()
		return
	}
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		c.Abort()
		return
	}
	var input struct {
		NewAdminUsername string `json:"new_admin_username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	tempGame, err := model.GetGameByID(uint(gameID))
	tempGame.ID = 0
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find game",
		}})
		c.Abort()
		return
	}
	if input.NewAdminUsername == "" {
		err = tempGame.Create([]uint{userID})
	} else {
		newAdmin, err := model.GetUserByUsername(input.NewAdminUsername)
		if err != nil {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "cannot find user",
			}})
			c.Abort()
			return
		}
		err = tempGame.Create([]uint{newAdmin.ID})
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create game"})
			c.Abort()
			return
		}
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create game"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"id": tempGame.ID})
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
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		gameData := gin.H{
			"base_contest": gin.H{
				"id":      game.ID,
				"game_id": game.BaseContest.GameID,
				"states": gin.H{
					"assign_ai_enabled":                  game.BaseContest.States.AssignAiEnabled,
					"commit_ai_enabled":                  game.BaseContest.States.CommitAiEnabled,
					"contest_script_environment_enabled": game.BaseContest.States.ContestScriptEnvironmentEnabled,
					"private_match_enabled":              game.BaseContest.States.PrivateMatchEnabled,
					"public_match_enabled":               game.BaseContest.States.PublicMatchEnabled,
					"test_match_enabled":                 game.BaseContest.States.TestMatchEnabled,
				},
			},
			"id":       game.ID,
			"metadata": basecontest.ConvertStruct(game.Metadata),
		}
		gamesList = append(gamesList, gameData)
	}

	c.JSON(200, gamesList)
}

func getGameSettings(c *gin.Context) {
	game.RetGameSettings(c)
}

func updateGameLogic(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	var input struct {
		BuildGameLogicDockerfile string                `form:"build_game_logic_dockerfile"`
		RunGameLogicDockerfile   string                `form:"run_game_logic_dockerfile"`
		GameLogicFile            *multipart.FileHeader `form:"game_logic_file"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// 保存上传的文件
	filePath := fmt.Sprintf("/var/hiper/games/%d/game_logic/gamelogic.zip", gameID)
	if err := c.SaveUploadedFile(input.GameLogicFile, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err := model.UpdateGameByID(gameID, map[string]interface{}{
		"game_logic_build_dockerfile": input.BuildGameLogicDockerfile,
		"game_logic_run_dockerfile":   input.RunGameLogicDockerfile,
	})
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update game logic",
		}})
		c.Abort()
		return
	}
	game.RetGameSettings(c)
}

func updateMatchDetail(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	var input struct {
		Template string `json:"template"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
	}
	err := model.UpdateGameByID(gameID, map[string]interface{}{
		"match_detail_template": input.Template,
	})
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot update match detail",
		}})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()
}
