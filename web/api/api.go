package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ApiListenHttp starts the HTTP server
func ApiListenHttp() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:    []string{"Content-Type", "Access-Token", "Authorization"},
		MaxAge:          6 * time.Hour,
	}))

	addUserRoutes(r)
	addPermissionRoutes(r)
	addGameRoutes(r)
	addContestRoutes(r)
	addBaseContestRoutes(r)

	r.Run(":8080")
}

func ApiListenHttps() {
	r := gin.Default()

	r.Use(cors.Default())
	addUserRoutes(r)
	addPermissionRoutes(r)
	addGameRoutes(r)
	addContestRoutes(r)
	addBaseContestRoutes(r)

	//请使用开发文档中的web-获取ssl证书来获得证书并将本地路径填入config中
	r.RunTLS(":8080", viper.GetString("net.certpath"), viper.GetString("net.keypath"))
}

func addUserRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/user/request-verification-code", requestVerificationCode)
		v1.POST("/users", registerUser)
		v1.POST("/user/reset-email", resetEmail)
		v1.POST("/user/reset-password", resetPassword)
		v1.POST("/user/login", login)

		//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
		auth := v1.Group("/", loginVerify())
		{
			auth.GET("/users", searchUsers)
			auth.DELETE("/user/logout", logout)
			auth.GET("/users/:username", getTheUser)
			auth.GET("/user", getCurrentUser)
			auth.PATCH("/user", updateCurrentUser)
		}
	}
}

func addPermissionRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
		auth := v1.Group("/", loginVerify())
		{
			auth.PUT("/permissions/create_game_or_contest/:user_id", grantCreationPermission)
			auth.DELETE("/permissions/create_game_or_contest/:user_id", revokeCreationPermission)

		}
	}
}

func addGameRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/games", getGames)
		//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
		auth := v1.Group("/", loginVerify())
		{
			auth.POST("/games", createGame)
			auth.POST("/games/:id/fork", forkGame)
			//此后的路由都需要验证是否是管理员.在其内部，我们可以使用gameID := c.MustGet("gameID").(int)来获取当前游戏的ID
			auth = auth.Group("/", privilegeCheck())
			{
				auth.GET("/games/:id/settings", getGameSettings)
				auth.PATCH("/games/:id/game_logic", updateGameLogic)
				auth.PATCH("/games/:id/match_detail", updateMatchDetail)
			}
		}
	}
}

func addContestRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/contests", getContests)
		//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
		auth := v1.Group("/", loginVerify())
		{
			auth.POST("/contests", createContest)
			auth.PUT("/contests/:id/register", registerContest)
			auth.DELETE("/contests/:id/register", exitContest)
			auth = auth.Group("/", privilegeCheck())
			{
				auth.GET("/contests/:id/settings", getContestSettings)
				auth.PUT("/contests/:id/password", updateContestPassword)
			}
		}
	}
}

func addBaseContestRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		//此后的路由都需要验证是否登录.在其内部，我们可以使用userID := c.MustGet("userID").(int)来获取当前登录用户的ID
		auth := v1.Group("/", loginVerify())
		{
			auth.GET("/games/:id", getTheGame)//no
			auth.GET("/games/:id/ais", getAis)//ok
			auth.POST("/games/:id/ais", commitAi)//no
			auth.GET("/games/:id/ais/:ai_id", getTheAI)//ok
			auth.GET("/games/:id/ais/:ai_id/file", downloadTheAI)//pending
			auth.PUT("/games/:id/ais/:ai_id/note", editAiNote)//ok
			auth.GET("/games/:id/contestants", getContestants)//no
			auth.PUT("/games/:id/contestant/assigned_ai", assignAi)//no
			auth.GET("/games/:id/contestant", getCurrentContestant)//ok
			auth.DELETE("/games/:id/contestant/assigned_ai", revokeAssignedAi)//pending
			auth.GET("/games/:id/matches", getMatches)//no
			auth.GET("/games/:id/matches/:match_id", getMatch)//no
			auth.GET("/games/:id/sdks", getSdks)//ok
			//此后的路由都需要验证是否是管理员.在其内部，我们可以使用gameID := c.MustGet("gameID").(int)来获取当前游戏的ID
			auth = auth.Group("/", privilegeCheck())
			{
				auth.DELETE("/games/:id", deleteGame)
				auth.POST("/games/:id/admins", addAdmin)
				auth.DELETE("/games/:id/admin", relinquishAdmin)
				auth.PUT("/games/:id/contest_script", updateGameScript)
				auth.PATCH("/games/:id/metadata", updateGameMetadata)
				auth.POST("/games/:id/sdks", addSdk)
				auth.GET("/games/:id/sdks/:sdk_id", getSdk)
				auth.DELETE("/games/:id/sdks/:sdk_id", deleteSdk)
				auth.PATCH("/games/:id/sdks/:sdk_id", updateSdk)
				auth.PATCH("/games/:id/states", updateGameStates)
			}
		}
	}
}
