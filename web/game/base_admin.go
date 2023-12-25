package game

import (
	"errors"
	"hiper-backend/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ErrorFor422 struct {
	Code   ErrorCodeFor422 `json:"code"`
	Field  string          `json:"field"`
	Detail string          `json:"detail"`
}

type ErrorCodeFor422 string

const (
	MissingField  ErrorCodeFor422 = "missing_field"
	Invalid       ErrorCodeFor422 = "invalid"
	AlreadyExists ErrorCodeFor422 = "already_exists"
)

func RetGameSettings(c *gin.Context) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	game, err := model.GetGameByID(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find game",
		}})
		c.Abort()
		return
	}
	baseContest, err := model.GetBaseContestByID(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find game",
		}})
		c.Abort()
		return
	}
	adminGame, _ := game.GetAdmins()
	admins := make([]map[string]interface{}, len(adminGame))
	for i, admin := range adminGame {
		admins[i] = map[string]interface{}{
			"avatar_url": admin.AvatarURL,
			"nickname":   admin.Nickname,
			"username":   admin.Username,
			"bio":        admin.Bio,
			"department": admin.Department,
			"name":       admin.Name,
			"permissions": map[string]bool{
				"can_create_game_or_contest": admin.Permissions.CanCreateGameOrContest,
			},
			"school": admin.School,
		}
	}
	rawSdk, err := baseContest.GetSdks()
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find sdks",
		}})
		c.Abort()
		return
	}
	sdks := make([]map[string]interface{}, len(rawSdk))
	for i, sdk := range rawSdk {
		sdks[i] = map[string]interface{}{
			"id":     sdk.ID,
			"name":   sdk.Name,
			"readme": sdk.Readme,
			"build_ai": map[string]interface{}{
				"status": sdk.BuildAi.Status,
			},
			"run_ai": map[string]interface{}{
				"status": sdk.RunAi.Status,
			},
		}
	}
	c.JSON(200, gin.H{
		"game_id": game.ID,
		"metadata": map[string]interface{}{
			"cover_url": game.Metadata.CoverUrl,
			"readme":    game.Metadata.Readme,
			"title":     game.Metadata.Title,
		},
		"states": map[string]interface{}{
			"assign_ai_enabled":                  game.BaseContest.States.AssignAiEnabled,
			"commit_ai_enabled":                  game.BaseContest.States.CommitAiEnabled,
			"contest_script_environment_enabled": game.BaseContest.States.ContestScriptEnvironmentEnabled,
			"private_match_enabled":              game.BaseContest.States.PrivateMatchEnabled,
			"public_match_enabled":               game.BaseContest.States.PublicMatchEnabled,
			"test_match_enabled":                 game.BaseContest.States.TestMatchEnabled,
		},
		"admins": admins,
		"contest_assets": map[string]interface{}{
			"contest_script": game.BaseContest.Script,
			"sdks":           sdks,
		},
		"id":           game.ID,
		"my_privilege": "admin",
		"game_assets": map[string]interface{}{
			"game_logic": map[string]interface{}{
				"build_game_logic": map[string]interface{}{
					"dockerfile": game.GameLogic.Build.Dockerfile,
					"status": map[string]interface{}{
						"state": game.GameLogic.Build.Status.State,
						"msg":   game.GameLogic.Build.Status.Msg,
					},
				},
				"run_game_logic": map[string]interface{}{
					"dockerfile": game.GameLogic.Build.Dockerfile,
					"status": map[string]interface{}{
						"state": game.GameLogic.Run.Status.State,
						"msg":   game.GameLogic.Run.Status.Msg,
					},
				},
			},
			"match_detail": map[string]interface{}{
				"template": game.MatchDetail.Template,
			},
		},
	})
	c.Abort()
}

func GetSdksFromKnownGame(c *gin.Context) (model.Sdk, error) {
	ingameID := c.MustGet("gameID").(int)
	gameID := uint(ingameID)
	str := c.Param("sdk_id")
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{})
	}
	sdkID := uint(num)
	baseContest, err := model.GetBaseContestByID(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find game",
		}})
		c.Abort()
		return model.Sdk{}, errors.New("cannot find game")
	}
	rawSdk, err := baseContest.GetSdks()
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find sdks",
		}})
		c.Abort()
		return model.Sdk{}, errors.New("cannot find sdks")
	}
	for _, sdk := range rawSdk {
		if sdk.ID == sdkID {
			return sdk, nil
		}
	}
	c.JSON(422, gin.H{"error": ErrorFor422{
		Code:  Invalid,
		Field: "cannot find sdk",
	}})
	c.Abort()
	return model.Sdk{}, errors.New("cannot find sdk")
}
