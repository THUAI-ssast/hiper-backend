package api

import (
	"hiper-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiListenHttp starts the HTTP server
func ApiListenHttp() {

	r := gin.Default()

	// TODO: routes
	r.POST("/users/request-verification-code", func(c *gin.Context) {
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

	r.POST("/users", func(c *gin.Context) {
		var jsonGetR GetRegister
		if err := c.ShouldBindJSON(&jsonGetR); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := jsonGetR.Email
		code := jsonGetR.Code
		password := jsonGetR.Password
		username := jsonGetR.Username
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
		} else if email_exist_codegiven(email) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "already_exists",
						"field": "email",
					},
				},
			})
		} else {
			userId := set_user(email, password, username)
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

	r.POST("/users/reset-email", func(c *gin.Context) {
		var jsonGetRE GetResetEmail
		if err := c.ShouldBindJSON(&jsonGetRE); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := jsonGetRE.Email
		code := jsonGetRE.Code
		newEmail := jsonGetRE.NewEmail
		if !email_exist_codegiven(email) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "email",
					},
				},
			})
		} else if !verify_email(newEmail) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "new_email",
					},
				},
			})
		} else if !code_match(code, email) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "verification_code",
					},
				},
			})
		} else {
			update_email(email, newEmail)
			c.JSON(200, gin.H{})
		}
	})

	r.POST("/users/reset-password", func(c *gin.Context) {
		var jsonGetRP GetResetPassword
		if err := c.ShouldBindJSON(&jsonGetRP); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := jsonGetRP.Email
		code := jsonGetRP.Code
		password := jsonGetRP.Password
		if !email_exist_codegiven(email) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "email",
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
		} else if !code_match(code, email) {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "invalid",
						"field": "verification_code",
					},
				},
			})
		} else {
			update_password(email, password)
			c.JSON(200, gin.H{})
		}
	})

	r.POST("/users/login", func(c *gin.Context) {
		var jsonGetL GetLogin
		if err := c.ShouldBindJSON(&jsonGetL); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := jsonGetL.Email
		username := jsonGetL.Username
		password := jsonGetL.Password
		if username == "" && email == "" {
			c.JSON(422, gin.H{
				"errors": []gin.H{
					{
						"code":  "missing_field",
						"field": "email and username",
					},
				},
			})
		} else if username != "" {
			if !username_exist_user(username) {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid",
							"field": "username",
						},
					},
				})
			} else if !password_match_username(username, password) {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid",
							"field": "password",
						},
					},
				})
			} else {
				userId := get_userId_username(username)
				token, _ := utils.GenToken((int64)(userId))
				c.JSON(200, gin.H{
					"access_token": token,
				})
			}
		} else {
			if !email_exist_user(email) {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid",
							"field": "email",
						},
					},
				})
			} else if !password_match_email(email, password) {
				c.JSON(422, gin.H{
					"errors": []gin.H{
						{
							"code":  "invalid",
							"field": "password",
						},
					},
				})
			} else {
				userId := get_userId_email(email)
				token, _ := utils.GenToken((int64)(userId))
				c.JSON(200, gin.H{
					"access_token": token,
				})
			}
		}
	})

	r.Run(":8080")
}
