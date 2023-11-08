package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiListenHttp starts the HTTP server
func ApiListenHttp() {

	r := gin.Default()

	// TODO: routes
	r.GET("/users/request-verification-code", func(c *gin.Context) {
		var jsonGetV GetRequestVerificationCode
		if err := c.ShouldBindJSON(&jsonGetV); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := jsonGetV.Email
		if verify_email(email) {
			send_email(email)
			c.JSON(200, gin.H{})
		} else {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "email",
					},
				},
			})
		}
	})
	r.GET("/users", func(c *gin.Context) {
		var jsonGetR GetRegister
		if err := c.ShouldBindJSON(&jsonGetR); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := jsonGetR.Email
		code := jsonGetR.Code
		password := jsonGetR.Password
		if !code_match(code, email) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "verification_code",
					},
				},
			})
		} else if !verify_password(password) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "password",
					},
				},
			})
		} else if !email_not_exist(email) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "already_exists",
						"field": "email",
					},
				},
			})
		} else {
			userId := set_user(email, password)
			if userId == -1 {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":   "missing_field",
							"field":  "email",
							"detail": "Delete from codegiven failed",
						},
					},
				})
			} else {
				c.JSON(200, gin.H{
					"user_id": userId,
				})
			}
		}

	})
	r.Run(":8080")
}
