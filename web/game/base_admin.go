package game

import (
	"hiper-backend/model"

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
	gameID := c.MustGet("gameID").(uint)
	game, err := model.GetGameById(gameID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find game",
		}})
		return
	}
	admins := make([]map[string]interface{}, len(game.Admins))
	for i, admin := range game.Admins {
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
	rawSdk, err := game.GetSdks()
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find sdks",
		}})
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
			"assign_ai_enabled":                  game.States.AssignAiEnabled,
			"commit_ai_enabled":                  game.States.CommitAiEnabled,
			"contest_script_environment_enabled": game.States.ContestScriptEnvironmentEnabled,
			"private_match_enabled":              game.States.PrivateMatchEnabled,
			"public_match_enabled":               game.States.PublicMatchEnabled,
			"test_match_enabled":                 game.States.TestMatchEnabled,
		},
		"admins": admins,
		"contest_assets": map[string]interface{}{
			"contest_script": game.Script,
			"sdks":           sdks,
		},
		"id":           game.ID,
		"my_privilege": "admin",
		"game_assets": map[string]interface{}{
			"game_logic": map[string]interface{}{
				"build_game_logic": map[string]interface{}{
					"status": map[string]interface{}{
						"state": game.GameLogic.Build.Status.State,
						"msg":   game.GameLogic.Build.Status.Msg,
					},
				},
				"run_game_logic": map[string]interface{}{
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
}
