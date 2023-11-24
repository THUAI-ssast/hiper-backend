package api

import (
	"github.com/gin-gonic/gin"
)

// ApiListenHttp starts the HTTP server
func ApiListenHttp() {
	r := gin.Default()

	addUserRoutes(r)
	// TODO: add more routes

	r.Run(":8080")
}

// addUserRoutes adds the routes for the user API
func addUserRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/user/request-verification-code", requestVerificationCode)
		v1.POST("/users", registerUser)
		// TODO: add more routes that don't require authentication

		auth := v1.Group("/", loginVerify())
		{
			auth.GET("/user", getCurrentUser)
			// TODO: add more routes that require authentication
		}
	}
}
