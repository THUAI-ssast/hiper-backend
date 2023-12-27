package contest

import (
	"github.com/gin-gonic/gin"

	"github.com/THUAI-ssast/hiper-backend/model"
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

func RetContestSettings(c *gin.Context) {
	incontestID := c.MustGet("gameID").(int)
	contestID := uint(incontestID)
	contest, err := model.GetContestByID(contestID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find contest",
		}})
		c.Abort()
		return
	}
	baseContest, err := model.GetBaseContestByID(contestID)
	if err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "cannot find contest",
		}})
		c.Abort()
		return
	}
	adminContest, _ := baseContest.GetAdmins()
	admins := make([]map[string]interface{}, len(adminContest))
	for i, admin := range adminContest {
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
		"base_contest": gin.H{
			"id":      contest.ID,
			"game_id": contest.ID,
			"states": map[string]interface{}{
				"assign_ai_enabled":                  contest.BaseContest.States.AssignAiEnabled,
				"commit_ai_enabled":                  contest.BaseContest.States.CommitAiEnabled,
				"contest_script_environment_enabled": contest.BaseContest.States.ContestScriptEnvironmentEnabled,
				"private_match_enabled":              contest.BaseContest.States.PrivateMatchEnabled,
				"public_match_enabled":               contest.BaseContest.States.PublicMatchEnabled,
				"test_match_enabled":                 contest.BaseContest.States.TestMatchEnabled,
			},
			"contest_assets": map[string]interface{}{
				"contest_script": contest.BaseContest.Script,
				"sdks":           sdks,
			},
		},
		"id": contest.ID,
		"metadata": map[string]interface{}{
			"cover_url": contest.Metadata.CoverUrl,
			"readme":    contest.Metadata.Readme,
			"title":     contest.Metadata.Title,
		},
		"admins":       admins,
		"my_privilege": "admin",
		"registration": map[string]interface{}{
			"registration_enabled": contest.Registration.RegistrationEnabled,
			"password":             contest.Registration.Password,
		},
	})
	c.Abort()
}
