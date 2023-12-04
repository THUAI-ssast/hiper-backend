package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ApiListenHttp starts the HTTP server
func ApiListenHttp() {
	r := gin.Default()

	addUserRoutes(r)
	addPermissionRoutes(r)
	addGameRoutes(r)
	// TODO: add more routes

	r.Run(":8080")
}

func ApiListenHttps() {
	r := gin.Default()

	addUserRoutes(r)
	addPermissionRoutes(r)
	addGameRoutes(r)
	// TODO: add more routes

	//请使用开发文档中的web-获取ssl证书来获得证书并将本地路径填入config中
	r.RunTLS(":8080", viper.GetString("net.certpath"), viper.GetString("net.keypath"))
}

func addUserRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/api/v1/user/request-verification-code", func(c *gin.Context) {
			requestVerificationCode(c)
		})

		v1.POST("/api/v1/users", func(c *gin.Context) {
			registerUser(c)
		})

		v1.POST("/api/v1/user/reset-email", func(c *gin.Context) {
			resetEmail(c)
		})

		v1.POST("/api/v1/user/reset-password", func(c *gin.Context) {
			resetPassword(c)
		})

		v1.POST("/api/v1/user/login", func(c *gin.Context) {
			login(c)
		})

		v1.GET("/api/v1/users/:username", func(c *gin.Context) {
			username := c.Param("username")
			getTheUser(c, username)
		})

		//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
		auth := v1.Group("/", loginVerify())
		{
			auth.GET("/api/v1/users", func(c *gin.Context) {
				searchUsers(c)
			})

			auth.DELETE("/api/v1/user/logout", func(c *gin.Context) {
				logout(c)
			})

			auth.GET("/api/v1/user", func(c *gin.Context) {
				getCurrentUser(c)
			})

			auth.PATCH("/api/v1/user", func(c *gin.Context) {
				updateCurrentUser(c)
			})
		}
	}
}

func addPermissionRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
		auth := v1.Group("/", loginVerify())
		{
			auth.PUT("/api/v1/permissions/create_game_or_contest/:user_id", func(c *gin.Context) {
				author_id := c.Param("user_id")
				grantCreationPermission(c, author_id)
			})

			auth.DELETE("/api/v1/permissions/create_game_or_contest/:user_id", func(c *gin.Context) {
				author_id := c.Param("user_id")
				revokeCreationPermission(c, author_id)
			})
		}
	}
}

func addGameRoutes(r *gin.Engine) {
	// auth.POST("/api/v1/games", func(c *gin.Context) {
	// 	create_game(c)
	// })

	// auth.POST("/api/v1/games/:id/fork", func(c *gin.Context) {
	// 	game_id := c.Param("id")
	// 	fork_game(c, game_id)
	// })
}
