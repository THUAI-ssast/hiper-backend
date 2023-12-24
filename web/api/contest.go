package api

import (
	"hiper-backend/model"

	"github.com/gin-gonic/gin"
)

//TODO:create contest要import"hiper-backend/mq"，在发送200前要先mq.SendBuildGameMsg(model.Ctx, tempGame.ID)，参见creategame最后几句

func getContests(c *gin.Context) {
	contests, err := model.GetContests()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	var contestsList []gin.H
	for _, contest := range contests {
		userID := c.MustGet("userID").(int)
		pri, err := contest.GetPrivilege(uint(userID))
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
			"id":           contest.ID,
			"metadata":     contest.Metadata,
			"my_privilege": pri,
		}
		contestsList = append(contestsList, contestData)
	}

	c.JSON(200, contestsList)
}
