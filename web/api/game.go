package api

import (
	"hiper-backend/game"
	"hiper-backend/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createGame(c *gin.Context) {
	inuserID := c.MustGet("userID").(int)
	userID := uint(inuserID)
	userr, err := model.GetUserById(userID)
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
		err = model.CreateGame(&tempGame, []uint{userID})
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
		err = model.CreateGame(&tempGame, []uint{newAdmin.ID})
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
	userr, err := model.GetUserById(userID)
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
	tempGame, err := model.GetGameById(uint(gameID))
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
		err = model.CreateGame(&tempGame, []uint{userID})
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
		err = model.CreateGame(&tempGame, []uint{newAdmin.ID})
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

func getGameSettings(c *gin.Context) {
	game.RetGameSettings(c)
}

func updateGameLogic(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	var input struct {
		BuildGameLogicDockerfile string `json:"build_game_logic_dockerfile"`
		RunGameLogicDockerfile   string `json:"run_game_logic_dockerfile"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	err := model.UpdateGameById(gameID, map[string]interface{}{
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
	err := model.UpdateGameById(gameID, map[string]interface{}{
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
