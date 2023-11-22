package api

import (
	"hiper-backend/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ApiListenHttp starts the HTTP server
func ApiListenHttp() {
	r := gin.Default()

	// TODO: routes
	r.POST("/api/v1/user/request-verification-code", func(c *gin.Context) {
		request_verification_code(c)
	})

	r.POST("/api/v1/users", func(c *gin.Context) {
		register_user(c)
	})

	r.POST("/api/v1/user/reset-email", func(c *gin.Context) {
		reset_email(c)
	})

	r.POST("/api/v1/user/reset-password", func(c *gin.Context) {
		reset_password(c)
	})

	r.POST("/api/v1/user/login", func(c *gin.Context) {
		login(c)
	})

	r.GET("/api/v1/users/:username", func(c *gin.Context) {
		username := c.Param("username")
		get_the_user(c, username)
	})

	//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
	authorized := r.Group("/")
	authorized.Use(model.Login_verify())
	{
		authorized.GET("/api/v1/users", func(c *gin.Context) {
			search_users(c)
		})

		authorized.DELETE("/api/v1/user/logout", func(c *gin.Context) {
			logout(c)
		})

		authorized.GET("/api/v1/user", func(c *gin.Context) {
			get_current_user(c)
		})

		authorized.PATCH("/api/v1/user", func(c *gin.Context) {
			update_current_user(c)
		})

		authorized.PUT("/api/v1/permissions/create_game_or_contest/:user_id", func(c *gin.Context) {
			author_id := c.Param("user_id")
			grant_creation_permission(c, author_id)
		})

		authorized.DELETE("/api/v1/permissions/create_game_or_contest/:user_id", func(c *gin.Context) {
			author_id := c.Param("user_id")
			revoke_creation_permission(c, author_id)
		})

		authorized.POST("/api/v1/games", func(c *gin.Context) {
			create_game(c)
		})

		authorized.POST("/api/v1/games/:id/fork", func(c *gin.Context) {
			game_id := c.Param("id")
			fork_game(c, game_id)
		})
	}
	//HTTP server
	//r.Run(":8080")
	//HTTPS server
	//请使用开发文档中的web-获取ssl证书来获得证书并将本地路径填入config中
	r.RunTLS(":8080", viper.GetString("net.certpath"), viper.GetString("net.keypath"))
}
