package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hiper-backend/mail"
	"hiper-backend/model"
	"hiper-backend/user"
)

func requestVerificationCode(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	email := input.Email
	if !mail.IsValidEmail(email) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "email",
		}})
		c.Abort()
		return
	}

	if err := user.SendVerificationCode(email); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}

func registerUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Code     string `json:"verification_code" binding:"required"`
		Password string `json:"password" binding:"required"`
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	email := input.Email
	code := input.Code
	password := input.Password
	username := input.Username

	if !user.IsCodeMatch(code, email) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "verification_code",
		}})
		c.Abort()
		return
	}
	if !mail.IsValidEmail(email) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "email",
		}})
		c.Abort()
		return
	}
	if !user.IsValidPassword(password) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "password",
		}})
		c.Abort()
		return
	}
	if _, err := model.GetUserByEmail(email); err == nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  AlreadyExists,
			Field: "email",
		}})
		c.Abort()
		return
	}
	if _, err := model.GetUserByUsername(username); err == nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  AlreadyExists,
			Field: "username",
		}})
		c.Abort()
		return
	}

	if _, err := user.RegisterUser(username, email, password); err != nil {
		c.JSON(500, gin.H{"error": "Failed to register user"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"username": username,
	})
	c.Abort()

}

func resetEmail(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Code     string `json:"verification_code" binding:"required"`
		NewEmail string `json:"new_email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	email := input.Email
	code := input.Code
	newEmail := input.NewEmail
	if email == newEmail {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:   Invalid,
			Field:  "new_email",
			Detail: "new email is the same as the old one",
		}})
		c.Abort()
		return
	}
	if _, err := model.GetUserByEmail(email); err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "email",
		}})
		c.Abort()
		return
	}
	if !mail.IsValidEmail(newEmail) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "new_email",
		}})
		c.Abort()
		return
	}
	if !user.IsCodeMatch(code, email) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "verification_code",
		}})
		c.Abort()
		return
	}
	if _, err := model.GetUserByEmail(newEmail); err == nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  AlreadyExists,
			Field: "new_email",
		}})
		c.Abort()
		return
	}

	if model.UpdateUserByEmail(email, map[string]interface{}{"email": newEmail}) != nil {
		c.JSON(500, gin.H{"error": "Failed to update user info"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()

}

func resetPassword(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Code     string `json:"verification_code" binding:"required"`
		Password string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	email := input.Email
	code := input.Code
	password := input.Password
	if _, err := model.GetUserByEmail(email); err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "email",
		}})
		c.Abort()
		return
	}
	if !user.IsValidPassword(password) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "new_password",
		}})
		c.Abort()
		return
	}
	if !user.IsCodeMatch(code, email) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "verification_code",
		}})
		c.Abort()
		return
	}

	if model.UpdateUserByEmail(email, map[string]interface{}{"password": user.HashPassword(password)}) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
	c.Abort()
}

func login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	email := input.Email
	username := input.Username
	password := input.Password

	if username == "" && email == "" {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  MissingField,
			Field: "email or username",
		}})
		c.Abort()
		return
	}

	if username != "" {
		// login by username
		var userID uint
		var valid bool
		if userID, valid = user.CheckPasswordByUsername(username, password); !valid {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "password",
			}})
			c.Abort()
			return
		}
		token, _ := GenToken((int64)(userID))
		c.JSON(200, gin.H{
			"access_token": token,
		})
		c.Abort()
		return
	} else {
		// login by email
		var userID uint
		var valid bool
		if userID, valid = user.CheckPasswordByEmail(email, password); !valid {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "password",
			}})
			c.Abort()
			return
		}
		token, _ := GenToken((int64)(userID))
		c.JSON(200, gin.H{
			"access_token": token,
		})
		c.Abort()
	}
}

func logout(c *gin.Context) {
	// nothing to do now
	c.JSON(200, gin.H{})
	c.Abort()
}

func searchUsers(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID := 0
	var err error
	if userIDStr != "" {
		userID, err = strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}
	}
	email := c.Query("email")
	username := c.Query("username")
	var users []model.User
	var answer []map[string]interface{}
	if email == "" && username == "" && userID == 0 {
		users, err = model.SearchUsers("", []string{"email"})
		if err != nil {
			c.JSON(404, gin.H{})
			c.Abort()
			return
		}
		for _, usr := range users {
			answer = append(answer, map[string]interface{}{
				"avatar_url": usr.AvatarURL,
				"username":   usr.Username,
				"bio":        usr.Bio,
				"department": usr.Department,
				"name":       usr.Name,
				"permissions": map[string]interface{}{
					"can_create_game_or_contest": usr.Permissions.CanCreateGameOrContest,
				},
				"school": usr.School,
			})
		}
		c.JSON(200, gin.H{
			"answer": answer,
		})
		c.Abort()
	} else if userID != 0 {
		usr, err := model.GetUserByID((uint)(userID))
		if err != nil {
			c.JSON(404, gin.H{})
			c.Abort()
			return
		}
		users = append(users, usr)
		for _, usr := range users {
			answer = append(answer, map[string]interface{}{
				"avatar_url": usr.AvatarURL,
				"username":   usr.Username,
				"bio":        usr.Bio,
				"department": usr.Department,
				"name":       usr.Name,
				"permissions": map[string]interface{}{
					"can_create_game_or_contest": usr.Permissions.CanCreateGameOrContest,
				},
				"school": usr.School,
			})
		}
		c.JSON(200, gin.H{
			"answer": answer,
		})
		c.Abort()
	} else if username != "" {
		users, err = model.SearchUsers(username, []string{"username"})
		if err != nil {
			c.JSON(404, gin.H{})
			c.Abort()
			return
		}
		for _, usr := range users {
			answer = append(answer, map[string]interface{}{
				"avatar_url": usr.AvatarURL,
				"username":   usr.Username,
				"bio":        usr.Bio,
				"department": usr.Department,
				"name":       usr.Name,
				"permissions": map[string]interface{}{
					"can_create_game_or_contest": usr.Permissions.CanCreateGameOrContest,
				},
				"school": usr.School,
			})
		}
		c.JSON(200, gin.H{
			"answer": answer,
		})
		c.Abort()
	} else {
		users, err = model.SearchUsers(email, []string{"email"})
		if err != nil {
			c.JSON(404, gin.H{})
			c.Abort()
			return
		}
		for _, usr := range users {
			answer = append(answer, map[string]interface{}{
				"avatar_url": usr.AvatarURL,
				"username":   usr.Username,
				"bio":        usr.Bio,
				"department": usr.Department,
				"name":       usr.Name,
				"permissions": map[string]interface{}{
					"can_create_game_or_contest": usr.Permissions.CanCreateGameOrContest,
				},
				"school": usr.School,
			})
		}
		c.JSON(200, gin.H{
			"answer": answer,
		})
		c.Abort()
	}
}

func getTheUser(c *gin.Context) {
	username := c.Param("username")
	usr, err := model.GetUserByUsername(username)
	if err != nil {
		c.JSON(404, gin.H{})
		c.Abort()
		return
	} else {
		contestant, err := model.GetContestantsByUserId(usr.ID)
		if err != nil {
			c.JSON(404, gin.H{})
			c.Abort()
			return
		}
		registered := make([]map[string]interface{}, 0)
		for _, ct := range contestant {
			if ct.ContestId == 0 {
				continue
			}
			game, err := model.GetContestById(ct.ContestId)
			if err != nil {
				c.JSON(404, gin.H{})
				c.Abort()
				return
			}
			myPrivilege := "registered"
			pri, _ := model.GetContestPrivilege(ct.ContestId, usr.ID)
			if pri == "admin" {
				myPrivilege = "admin"
			}
			registered = append(registered, map[string]interface{}{
				"game_id": ct.ContestId,
				"metadata": map[string]interface{}{
					"cover_url": game.Metadata.CoverUrl,
					"readme":    game.Metadata.Readme,
					"title":     game.Metadata.Title,
				},
				"states": map[string]interface{}{
					"commit_ai_enabled":                  game.States.CommitAiEnabled,
					"assign_ai_enabled":                  game.States.AssignAiEnabled,
					"public_match_enabled":               game.States.PublicMatchEnabled,
					"contest_script_environment_enabled": game.States.ContestScriptEnvironmentEnabled,
					"private_match_enabled":              game.States.PrivateMatchEnabled,
					"test_match_enabled":                 game.States.TestMatchEnabled,
				},
				"id":           ct.ID,
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
			"username":            usr.Username,
			"email":               usr.Email,
			"contests_registered": registered,
		})
		c.Abort()
	}
}

func getCurrentUser(c *gin.Context) {
	userIDs, _ := c.Get("userID")
	userID, _ := userIDs.(int)
	usr, err := model.GetUserByID((uint)(userID))
	if err != nil {
		c.JSON(404, gin.H{})
		c.Abort()
		return
	} else {
		contestant, err := model.GetContestantsByUserId(usr.ID)
		if err != nil {
			c.JSON(404, gin.H{})
			c.Abort()
			return
		}
		registered := make([]map[string]interface{}, 0)
		for _, ct := range contestant {
			if ct.ContestId == 0 {
				continue
			}
			game, err := model.GetContestById(ct.ContestId)
			if err != nil {
				c.JSON(404, gin.H{})
				c.Abort()
				return
			}
			myPrivilege := "registered"
			pri, _ := model.GetContestPrivilege(ct.ContestId, usr.ID)
			if pri == "admin" {
				myPrivilege = "admin"
			}
			registered = append(registered, map[string]interface{}{
				"game_id": ct.ContestId,
				"metadata": map[string]interface{}{
					"cover_url": game.Metadata.CoverUrl,
					"readme":    game.Metadata.Readme,
					"title":     game.Metadata.Title,
				},
				"states": map[string]interface{}{
					"commit_ai_enabled":                  game.States.CommitAiEnabled,
					"assign_ai_enabled":                  game.States.AssignAiEnabled,
					"public_match_enabled":               game.States.PublicMatchEnabled,
					"contest_script_environment_enabled": game.States.ContestScriptEnvironmentEnabled,
					"private_match_enabled":              game.States.PrivateMatchEnabled,
					"test_match_enabled":                 game.States.TestMatchEnabled,
				},
				"id":           ct.ID,
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
			"username":            usr.Username,
			"email":               usr.Email,
			"contests_registered": registered,
		})
		c.Abort()
	}
}

func updateCurrentUser(c *gin.Context) {
	userIDs, _ := c.Get("userID")
	userID, _ := userIDs.(int)
	_, err := model.GetUserByID((uint)(userID))
	if err != nil {
		c.JSON(404, gin.H{})
		c.Abort()
		return
	} else {
		var input struct {
			Avatar_url string `json:"avatar_url"`
			Bio        string `json:"bio"`
			Department string `json:"department"`
			Name       string `json:"name"`
			School     string `json:"school"`
			Username   string `json:"username"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "json",
			}})
			c.Abort()
			return
		}

		if usr, err := model.GetUserByUsername(input.Username, "ID"); err == nil && usr.ID != (uint)(userID) {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  AlreadyExists,
				Field: "username",
			}})
			c.Abort()
			return
		}

		if !user.IsValidURL(input.Avatar_url) {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "avatar_url",
			}})
			c.Abort()
			return
		}

		updates := map[string]interface{}{
			"avatar_url": input.Avatar_url,
			"username":   input.Username,
			"bio":        input.Bio,
			"department": input.Department,
			"name":       input.Name,
			"school":     input.School,
		}

		for key, value := range updates {
			if len(value.(string)) > 100 { // assuming 100 is the maximum length
				c.JSON(422, gin.H{"error": ErrorFor422{
					Code:  Invalid,
					Field: key,
				}})
				c.Abort()
				return
			}
		}

		if len(updates) > 0 {
			err = model.UpdateUserByID((uint)(userID), updates)
			if err != nil {
				c.JSON(422, gin.H{"error": ErrorFor422{
					Code:  Invalid,
					Field: "update failed",
				}})
				c.Abort()
				return
			}
		}
		usr, _ := model.GetUserByID((uint)(userID))
		c.JSON(200, gin.H{
			"avatar_url": usr.AvatarURL,
			"username":   usr.Username,
			"bio":        usr.Bio,
			"department": usr.Department,
			"name":       usr.Name,
			"permissions": gin.H{
				"can_create_game_or_contest": usr.Permissions.CanCreateGameOrContest,
			},
			"school":              usr.School,
			"contests_registered": "", //usr.ContestsRegistered,
			"email":               usr.Email,
		})
		c.Abort()
	}
}
