package api

import (
	"hiper-backend/contest"
	"hiper-backend/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createContest(c *gin.Context) {
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
		GameID           uint   `json:"game_id"`
		NewAdminUsername string `json:"new_admin_username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	tempContest := model.Contest{}
	if input.NewAdminUsername == "" {
		err = tempContest.Create(input.GameID, []uint{userID})
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
		err = tempContest.Create(input.GameID, []uint{newAdmin.ID})
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create Contest"})
			c.Abort()
			return
		}
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create Contest"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"id": tempContest.ID})
	c.Abort()
}

func registerContest(c *gin.Context) {
	inuserID := c.MustGet("userID").(int)
	userID := uint(inuserID)
	usr, err := model.GetUserByID(userID)
	if err != nil {
		c.JSON(401, gin.H{})
		c.Abort()
		return
	}
	gameIDtem := c.Param("id")
	gameIDUint, _ := strconv.ParseUint(gameIDtem, 10, 32)
	gameID := uint(gameIDUint)

	var input struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	contest, err := model.GetContestByID(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:   Invalid,
			Field:  "game",
			Detail: "game not found",
		}})
		c.Abort()
		return
	}
	if !contest.Registration.RegistrationEnabled ||
		(contest.Registration.Password != input.Password && contest.Registration.Password != "") {
		c.JSON(401, gin.H{})
		c.Abort()
		return
	}

	contest.RegisteredUsers = append(contest.RegisteredUsers, usr)

	// 更新 Contest
	err = contest.Update(map[string]interface{}{
		"RegisteredUsers": contest.RegisteredUsers,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to register Contest"})
		c.Abort()
		return
	}

	// 更新 Contestant
	contestant := model.Contestant{
		BaseContestID: contest.ID,
		UserID:        usr.ID,
	}
	err = contestant.Create()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to register Contest"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
}

func exitContest(c *gin.Context) {
	inuserID := c.MustGet("userID").(int)
	userID := uint(inuserID)
	usr, err := model.GetUserByID(userID)
	if err != nil {
		c.JSON(401, gin.H{})
		c.Abort()
		return
	}
	gameIDtem := c.Param("id")
	gameIDUint, _ := strconv.ParseUint(gameIDtem, 10, 32)
	gameID := uint(gameIDUint)

	contest, err := model.GetContestByID(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:   Invalid,
			Field:  "game",
			Detail: "game not found",
		}})
		c.Abort()
		return
	}
	for i, user := range contest.RegisteredUsers {
		if user.ID == usr.ID {
			contest.RegisteredUsers = append(contest.RegisteredUsers[:i], contest.RegisteredUsers[i+1:]...)
			break
		}
	}
	// 更新 Contest
	err = contest.Update(map[string]interface{}{
		"RegisteredUsers": contest.RegisteredUsers,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to exit Contest"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"id": contest.ID})
}

func getContestSettings(c *gin.Context) {
	contest.RetContestSettings(c)
}

func updateContestPassword(c *gin.Context) {
	inuserID := c.MustGet("userID").(int)
	userID := uint(inuserID)
	usr, err := model.GetUserByID(userID)
	if err != nil {
		c.JSON(401, gin.H{})
		c.Abort()
		return
	}
	gameIDtem := c.Param("id")
	gameIDUint, _ := strconv.ParseUint(gameIDtem, 10, 32)
	gameID := uint(gameIDUint)

	var input struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	contest, err := model.GetContestByID(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:   Invalid,
			Field:  "game",
			Detail: "game not found",
		}})
		c.Abort()
		return
	}
	baseContest, err := model.GetBaseContestByID(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:   Invalid,
			Field:  "game",
			Detail: "game not found",
		}})
		c.Abort()
		return
	}
	admin, _ := baseContest.IsAdmin(usr.ID)
	if !admin {
		c.JSON(401, gin.H{})
		c.Abort()
		return
	}

	// 更新 Contest
	err = contest.Update(map[string]interface{}{
		"Password": input.Password,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to update Contest password"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"id": contest.ID})
}

func getContests(c *gin.Context) {
	contests, err := model.GetContests()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	var contestsList []gin.H
	for _, contest := range contests {
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		contestData := gin.H{
			"base_contest": gin.H{
				"id":      contest.ID,
				"game_id": contest.BaseContest.GameID,
				"states": gin.H{
					"assign_ai_enabled":                  contest.BaseContest.States.AssignAiEnabled,
					"commit_ai_enabled":                  contest.BaseContest.States.CommitAiEnabled,
					"contest_script_environment_enabled": contest.BaseContest.States.ContestScriptEnvironmentEnabled,
					"private_match_enabled":              contest.BaseContest.States.PrivateMatchEnabled,
					"public_match_enabled":               contest.BaseContest.States.PublicMatchEnabled,
					"test_match_enabled":                 contest.BaseContest.States.TestMatchEnabled,
				},
			},
			"id":       contest.ID,
			"metadata": contest.Metadata,
		}
		contestsList = append(contestsList, contestData)
	}

	c.JSON(200, contestsList)
}
