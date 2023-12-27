package user

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/THUAI-ssast/hiper-backend/web/mail"
	"github.com/THUAI-ssast/hiper-backend/web/model"
)

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendVerificationCode(email string) error {
	code := GenValidateCode(6)
	if err := model.SaveVerificationCode(code, email, 5); err != nil {
		return err
	}
	if err := mail.SendVerificationCode(code, email); err != nil {
		return err
	}
	return nil
}

// RegisterUser registers a user.
// It returns the user id and an error.
// The error is nil if the registration is successful.
func RegisterUser(username string, email string, password string) (uint, error) {
	hashedPassword := HashPassword(password)
	user := model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	if err := user.Create(); err != nil {
		return 0, err
	}
	return user.ID, nil
}

func ReturnWithUser(c *gin.Context, usr model.User, err error) {
	if err != nil {
		c.JSON(404, gin.H{})
		c.Abort()
		return
	} else {
		baseContests, err := usr.GetContestRegistered()
		if err != nil {
			c.JSON(404, gin.H{})
			c.Abort()
			return
		}
		registered := make([]map[string]interface{}, 0)
		for _, game := range baseContests {
			if err != nil {
				c.JSON(404, gin.H{})
				c.Abort()
				return
			}
			myPrivilege := "registered"
			pri, _ := game.GetPrivilege(usr.ID)
			if pri == "admin" {
				myPrivilege = "admin"
			}
			registered = append(registered, map[string]interface{}{
				"game_id": game.ID,
				"metadata": map[string]interface{}{
					"cover_url": game.Metadata.CoverUrl,
					"readme":    game.Metadata.Readme,
					"title":     game.Metadata.Title,
				},
				"states": map[string]interface{}{
					"commit_ai_enabled":                  game.BaseContest.States.CommitAiEnabled,
					"assign_ai_enabled":                  game.BaseContest.States.AssignAiEnabled,
					"public_match_enabled":               game.BaseContest.States.PublicMatchEnabled,
					"contest_script_environment_enabled": game.BaseContest.States.ContestScriptEnvironmentEnabled,
					"private_match_enabled":              game.BaseContest.States.PrivateMatchEnabled,
					"test_match_enabled":                 game.BaseContest.States.TestMatchEnabled,
				},
				"id":           game.ID,
				"my_privilege": myPrivilege,
			})
		}
		c.JSON(200, gin.H{
			"avatar_url": usr.AvatarURL,
			"bio":        usr.Bio,
			"department": usr.Department,
			"name":       usr.Name,
			"permissions": map[string]bool{
				"can_create_game_or_contest": usr.Permissions.CanCreateGameOrContest,
			},
			"school":              usr.School,
			"nickname":            usr.Nickname,
			"username":            usr.Username,
			"email":               usr.Email,
			"contests_registered": registered,
		})
		c.Abort()
	}
}
