package api

import (
	"github.com/gin-gonic/gin"
)

// ApiListenHttp starts the HTTP server
func ApiListenHttp() {

	r := gin.Default()

	// TODO: routes
	r.POST("/users/request-verification-code", func(c *gin.Context) {
		request_verification_code(c)
	})

	r.POST("/users", func(c *gin.Context) {
		register_user(c)
	})

	r.POST("/users/reset-email", func(c *gin.Context) {
		reset_email(c)
	})

	r.POST("/users/reset-password", func(c *gin.Context) {
		reset_password(c)
	})

	r.POST("/users/login", func(c *gin.Context) {
		login(c)
	})

	r.DELETE("/token", func(c *gin.Context) {
		logout(c)
	})

	r.Run(":8080")
}
