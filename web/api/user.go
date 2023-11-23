package api

import (
	"github.com/gin-gonic/gin"

	"hiper-backend/mail"
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
		c.JSON(422, gin.H{"errors": []ErrorFor422{{
			Code:  Invalid,
			Field: "email",
		}}})
		return
	}

	if err := user.SendVerificationCode(email); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{})
}

// TODO
func registerUser(c *gin.Context) {
}

// TODO
func getCurrentUser(c *gin.Context) {
}
