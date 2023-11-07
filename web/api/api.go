package api

import (
	"github.com/gin-gonic/gin"
)

// ApiListenHttp starts the HTTP server
func ApiListenHttp() {
	r := gin.Default()

	// TODO: routes

	r.Run(":8080")
}
