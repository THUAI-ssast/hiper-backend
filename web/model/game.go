package model

/*
import (
	"github.com/gin-gonic/gin"
)

func Contains(slice []interface{}, val interface{}) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func SendGameStatus(c *gin.Context, gameID int, userID int) {
	gameData, ok := SelectMySql("game", map[string]interface{}{"game_id": gameID})
	if !ok || len(gameData) == 0 {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "game_not_found",
				},
			},
		})
		return
	}

	var my_privilege string
	admins_id, ok := SelectAndParseJSONArray("game", map[string]interface{}{"game_id": gameID}, "admins_id")
	if !ok || len(admins_id) == 0 {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "admins_id_not_found",
				},
			},
		})
		return
	}
	if !Contains(admins_id, userID) {
		my_privilege = "registered"
	} else {
		my_privilege = "admin"
	}

	c.JSON(200, gin.H{
		"game_id": gameID,
		"metadata": gin.H{
			"cover_url": gameData[0]["cover_url"],
			"readme":    gameData[0]["readme"],
			"title":     gameData[0]["title"],
		},
		"states": gin.H{
			"assign_ai_enabled":                  gameData[0]["states_assign_ai_enabled"] == "1",
			"commit_ai_enabled":                  gameData[0]["states_commit_ai_enabled"] == "1",
			"contest_script_environment_enabled": gameData[0]["states_contest_script_environment_enabled"] == "1",
			"private_match_enabled":              gameData[0]["states_private_match_enabled"] == "1",
			"public_match_enabled":               gameData[0]["states_public_match_enabled"] == "1",
			"test_match_enabled":                 gameData[0]["states_test_match_enabled"] == "1",
		},
		"my_privilege": my_privilege,
		"id":           gameID,
	})
}
*/
