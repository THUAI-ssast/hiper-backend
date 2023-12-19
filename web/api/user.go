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
		return
	}

	email := input.Email
	if !mail.IsValidEmail(email) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "email",
		}})
		return
	}

	if err := user.SendVerificationCode(email); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
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
	}
	if !mail.IsValidEmail(email) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "email",
		}})
	}
	if !user.IsValidPassword(password) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "password",
		}})
	}
	if _, err := model.GetUserByEmail(email); err == nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  AlreadyExists,
			Field: "email",
		}})
	}
	if _, err := model.GetUserByUsername(username); err == nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  AlreadyExists,
			Field: "username",
		}})
	}

	u := model.User{
		Email:    email,
		Password: user.HashPassword(password),
		Username: username,
	}
	if err := model.CreateUser(&u); err != nil {
		c.JSON(500, gin.H{"error": "user creation failed"})
	}
	c.JSON(200, gin.H{
		"username": username,
	})

}

func resetEmail(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Code     string `json:"verification_code" binding:"required"`
		NewEmail string `json:"new_email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
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
	}
	if _, err := model.GetUserByEmail(email); err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "email",
		}})
	}
	if !mail.IsValidEmail(newEmail) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "new_email",
		}})
	}
	if !user.IsCodeMatch(code, email) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "verification_code",
		}})
	}
	if _, err := model.GetUserByEmail(newEmail); err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  AlreadyExists,
			Field: "new_email",
		}})
	}

	if model.UpdateUserByEmail(email, map[string]interface{}{"email": newEmail}) != nil {
		c.JSON(500, gin.H{"error": "Failed to update user info"})
	}
	c.JSON(200, gin.H{})

}

func resetPassword(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Code     string `json:"verification_code" binding:"required"`
		Password string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	email := input.Email
	code := input.Code
	password := input.Password
	if _, err := model.GetUserByEmail(email); err != nil {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "email",
		}})
	}
	if !user.IsValidPassword(password) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "new_password",
		}})
	}
	if !user.IsCodeMatch(code, email) {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  Invalid,
			Field: "verification_code",
		}})
	}

	if model.UpdateUserByEmail(email, map[string]interface{}{"password": user.HashPassword(password)}) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
	}
	c.JSON(200, gin.H{})
}

func login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	email := input.Email
	username := input.Username
	password := input.Password

	if username == "" && email == "" {
		c.JSON(422, gin.H{"error": ErrorFor422{
			Code:  MissingField,
			Field: "email or username",
		}})
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
		}
		token, _ := GenToken((int64)(userID))
		c.JSON(200, gin.H{
			"access_token": token,
		})
	} else {
		// login by email
		var userID uint
		var valid bool
		if userID, valid = user.CheckPasswordByEmail(email, password); !valid {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "password",
			}})
		}
		token, _ := GenToken((int64)(userID))
		c.JSON(200, gin.H{
			"access_token": token,
		})
	}
}

func logout(c *gin.Context) {
	// nothing to do now
	c.JSON(200, gin.H{})
}

func searchUsers(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID := 0
	var err error
	if userIDStr != "" {
		userID, err = strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
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
	} else if userID != 0 {
		usr, err := model.GetUserById((uint)(userID))
		if err != nil {
			c.JSON(404, gin.H{})
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
	} else if username != "" {
		users, err = model.SearchUsers(username, []string{"username"})
		if err != nil {
			c.JSON(404, gin.H{})
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
	} else {
		users, err = model.SearchUsers(email, []string{"email"})
		if err != nil {
			c.JSON(404, gin.H{})
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
	}
}

func getTheUser(c *gin.Context) {
	username := c.Param("username")
	usr, err := model.GetUserByUsername(username)
	if err != nil {
		c.JSON(404, gin.H{})
	} else {
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
			"contests_registered": "", //usr.ContestsRegistered,
		})
	}
}

func getCurrentUser(c *gin.Context) {
	userIDs, _ := c.Get("userID")
	userID, _ := userIDs.(int)
	usr, err := model.GetUserById((uint)(userID))
	if err != nil {
		c.JSON(404, gin.H{})
	} else {
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
			"contests_registered": "", //usr.ContestsRegistered,
		})
	}
}

func updateCurrentUser(c *gin.Context) {
	userIDs, _ := c.Get("userID")
	userID, _ := userIDs.(int)
	_, err := model.GetUserById((uint)(userID))
	if err != nil {
		c.JSON(404, gin.H{})
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
		}

		if usr, err := model.GetUserByUsername(input.Username, "ID"); err == nil && usr.ID != (uint)(userID) {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  AlreadyExists,
				Field: "username",
			}})
		}

		if !user.IsValidURL(input.Avatar_url) {
			c.JSON(422, gin.H{"error": ErrorFor422{
				Code:  Invalid,
				Field: "avatar_url",
			}})
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
			}
		}

		if len(updates) > 0 {
			err = model.UpdateUserById((uint)(userID), updates)
			if err != nil {
				c.JSON(422, gin.H{"error": ErrorFor422{
					Code:  Invalid,
					Field: "update failed",
				}})
			}
		}
		usr, _ := model.GetUserById((uint)(userID))
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
	}
}
