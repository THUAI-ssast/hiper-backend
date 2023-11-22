package api

import (
	"hiper-backend/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetCreateGame struct {
	New_username string `json:"new_admin_username"`
}

func create_game(c *gin.Context) {
	var userID int
	var err error
	userID = c.MustGet("userID").(int)
	if ret, valid := model.SelectMySql("user", map[string]interface{}{"user_id": userID}); !valid || ret[0]["authorization"] == "Regular user" {
		c.JSON(401, gin.H{})
		return
	}
	var jsongetCG GetCreateGame
	if err = c.ShouldBindJSON(&jsongetCG); err != nil {
		c.JSON(401, gin.H{})
		return
	}
	new_username := jsongetCG.New_username
	gameID, valid := model.InsertMySql("game", map[string]interface{}{})
	if !valid {
		c.JSON(401, gin.H{})
		return
	}
	if new_username != "" {
		if ret, valid := model.SelectMySql("user", map[string]interface{}{"username": new_username}); !valid {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "new_admin_username",
					},
				},
			})
			return
		} else {
			userID, err = strconv.Atoi(ret[0]["user_id"])
			if err != nil {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid",
							"field": "user_id",
						},
					},
				})
				return
			}
		}
	}
	if !model.AppendToJSONArrayInMySQL("user", "contests_id", gameID, map[string]interface{}{"user_id": userID}) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "append_game_to_user",
				},
			},
		})
		return
	} else {
		if !model.AppendToJSONArrayInMySQL("game", "admins_id", userID, map[string]interface{}{"game_id": gameID}) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "append_user_to_game",
					},
				},
			})
			return
		} else {
			model.SendGameStatus(c, gameID, c.MustGet("userID").(int))
		}
	}
}

func fork_game(c *gin.Context, game_ids string) {
	var userID int
	var err error
	userID = c.MustGet("userID").(int)
	if ret, valid := model.SelectMySql("user", map[string]interface{}{"user_id": userID}); !valid || ret[0]["authorization"] == "Regular user" {
		c.JSON(401, gin.H{})
		return
	}
	var jsongetCG GetCreateGame
	if err = c.ShouldBindJSON(&jsongetCG); err != nil {
		c.JSON(401, gin.H{})
		return
	}
	fork_game_id, err := strconv.Atoi(game_ids)
	if err != nil {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "game_id",
				},
			},
		})
		return
	}
	gameData, valid := model.SelectMySql("game", map[string]interface{}{"game_id": fork_game_id})
	if !valid || len(gameData) == 0 {
		c.JSON(404, gin.H{})
		return
	}
	new_username := jsongetCG.New_username
	if new_username != "" {
		if ret, valid := model.SelectMySql("user", map[string]interface{}{"username": new_username}); !valid {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "new_admin_username",
					},
				},
			})
			return
		} else {
			userID, err = strconv.Atoi(ret[0]["user_id"])
			if err != nil {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid",
							"field": "user_id",
						},
					},
				})
				return
			}
		}
	}

	gameID, valid := model.InsertMySql("game", map[string]interface{}{
		"metadata_cover_url":                        gameData[0]["cover_url"],
		"metadata_readme":                           gameData[0]["readme"],
		"metadata_title":                            gameData[0]["title"],
		"states_assign_ai_enabled":                  gameData[0]["states_assign_ai_enabled"],
		"states_commit_ai_enabled":                  gameData[0]["states_commit_ai_enabled"],
		"states_contest_script_environment_enabled": gameData[0]["states_contest_script_environment_enabled"],
		"states_private_match_enabled":              gameData[0]["states_private_match_enabled"],
		"states_public_match_enabled":               gameData[0]["states_public_match_enabled"],
		"states_test_match_enabled":                 gameData[0]["states_test_match_enabled"],
		"contest_assets_contest_script":             gameData[0]["contest_assets_contest_script"],
	})
	if !valid {
		c.JSON(401, gin.H{})
		return
	}

	sdkIDList, valid := model.SelectAndParseJSONArray("game", map[string]interface{}{"game_id": fork_game_id}, "contest_assets_sdks_id")
	if !valid {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "contest_assets_sdks_id",
				},
			},
		})
		return
	}
	//将sdk深度复制到新的game中
	for _, sdkID := range sdkIDList {
		sdkData, ok := model.SelectMySql("sdk", map[string]interface{}{"id": sdkID})
		if !ok || len(sdkData) == 0 {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "sdk_not_found",
					},
				},
			})
			return
		}
		newSdkID, valid := model.InsertMySql("sdk", map[string]interface{}{
			"name":           sdkData[0]["name"],
			"description":    sdkData[0]["description"],
			"build_ai_state": sdkData[0]["build_ai_state"],
			"run_ai_state":   sdkData[0]["run_ai_state"],
		})
		if !valid {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "sdk_creation_failed",
					},
				},
			})
			return
		}
		if !model.AppendToJSONArrayInMySQL("game", "contest_assets_sdks_id", newSdkID, map[string]interface{}{"game_id": gameID}) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "append_sdk_to_game",
					},
				},
			})
			return
		}
	}

	if !model.AppendToJSONArrayInMySQL("user", "contests_id", gameID, map[string]interface{}{"user_id": userID}) {
		c.JSON(422, gin.H{
			"errors": []gin.H{
				{
					"code":  "invalid",
					"field": "append_game_to_user",
				},
			},
		})
		return
	} else {
		if !model.AppendToJSONArrayInMySQL("game", "admins_id", userID, map[string]interface{}{"game_id": gameID}) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "append_user_to_game",
					},
				},
			})
			return
		} else {
			model.SendGameStatus(c, gameID, c.MustGet("userID").(int))
		}
	}
}
