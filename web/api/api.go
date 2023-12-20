package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ApiListenHttp starts the HTTP server
func ApiListenHttp() {
	r := gin.Default()

	r.Use(cors.Default())

	addUserRoutes(r)
	addPermissionRoutes(r)
	addGameRoutes(r)
	// TODO: add more routes

	r.Run(":8080")
}

func ApiListenHttps() {
	r := gin.Default()

	r.Use(cors.Default())
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
		v1.POST("/user/request-verification-code", func(c *gin.Context) {
			requestVerificationCode(c)
		})

		v1.POST("/users", func(c *gin.Context) {
			registerUser(c)
		})

		v1.POST("/user/reset-email", func(c *gin.Context) {
			resetEmail(c)
		})

		v1.POST("/user/reset-password", func(c *gin.Context) {
			resetPassword(c)
		})

		v1.POST("/user/login", func(c *gin.Context) {
			login(c)
		})

		//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
		auth := v1.Group("/", loginVerify())
		{
			auth.GET("/users", func(c *gin.Context) {
				searchUsers(c)
			})

			auth.DELETE("/user/logout", func(c *gin.Context) {
				logout(c)
			})

			auth.GET("/users/:username", func(c *gin.Context) {
				getTheUser(c)
			})

			auth.GET("/user", func(c *gin.Context) {
				getCurrentUser(c)
			})

			auth.PATCH("/user", func(c *gin.Context) {
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
			auth.PUT("/permissions/create_game_or_contest/:user_id", func(c *gin.Context) {
				grantCreationPermission(c)
			})

			auth.DELETE("/permissions/create_game_or_contest/:user_id", func(c *gin.Context) {
				revokeCreationPermission(c)
			})

		}
	}
}

func addGameRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
		auth := v1.Group("/", loginVerify())
		{
			auth.POST("/games", func(c *gin.Context) {
				createGame(c)
			})
			auth.POST("/games/:id/fork", func(c *gin.Context) {
				forkGame(c)
			})
			//此后的路由都需要验证是否是管理员.在其内部，我们可以使用gameID := c.MustGet("gameID").(int)来获取当前游戏的ID
			auth = auth.Group("/", privilegeCheck())
			{
				auth.DELETE("/games/:id", func(c *gin.Context) {
					deleteGame(c)
				})
				auth.POST("/games/:id/admins", func(c *gin.Context) {
					addAdmin(c)
				})
				auth.DELETE("/games/:id/admins", func(c *gin.Context) {
					relinquishAdmin(c)
				})
				auth.PUT("/games/:id/contest_script", func(c *gin.Context) {
					updateGameScript(c)
				})
				auth.PATCH("/games/:id/metadata", func(c *gin.Context) {
					updateGameMetadata(c)
				})
				auth.POST("/games/:id/sdks", func(c *gin.Context) {
					addSdk(c)
				})
				auth.GET("/games/:id/sdks/:sdk_id", func(c *gin.Context) {
					getSdk(c)
				})
				auth.DELETE("/games/:id/sdks/:sdk_id", func(c *gin.Context) {
					deleteSdk(c)
				})
				auth.PATCH("/games/:id/sdks/:sdk_id", func(c *gin.Context) {
					updateSdk(c)
				})
				auth.PATCH("/games/:id/states", func(c *gin.Context) {
					updateGameStates(c)
				})
				auth.GET("/games/:id/settings", func(c *gin.Context) {
					getGameSettings(c)
				})
				auth.PATCH("/games/:id/game_logic", func(c *gin.Context) {
					updateGameLogic(c)
				})
				auth.PATCH("/games/:id/match_detail", func(c *gin.Context) {
					updateMatchDetail(c)
				})
			}
		}
	}
}
