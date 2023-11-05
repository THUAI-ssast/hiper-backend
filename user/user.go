package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Main() {
	r := gin.Default()
	r.GET("/user", func(c *gin.Context) {
		name := c.DefaultQuery("name", "Hiper")
		c.String(http.StatusOK, fmt.Sprintf("Hello %s", name))
	})
	r.Run(":8000")
}
