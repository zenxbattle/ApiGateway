package router

import (
	"net/http"
	"xcode/clients"
	"xcode/configs"
	"xcode/controller"
	"xcode/middleware"
	"xcode/natsclient"

	"github.com/gin-gonic/gin"
	AuthUserAdminService "github.com/lijuuu/GlobalProtoXcode/AuthUserAdminService"

	ChallengeService "github.com/lijuuu/GlobalProtoXcode/ChallengeService"
	ProblemsService "github.com/lijuuu/GlobalProtoXcode/ProblemsService"
	"go.uber.org/zap"
)

func SetupRoutes(Router *gin.Engine, Clients *clients.ClientConnections, JWTSecret string, log *zap.Logger) {

	//Setup Client Instances
	NatsClient := natsclient.NewNatsClient(configs.LoadConfig().NATSURL, log)
	ProblemClient := ProblemsService.NewProblemsServiceClient(Clients.ConnProblem)
	UserClient := AuthUserAdminService.NewAuthUserAdminServiceClient(Clients.ConnUser)
	ChallengeClient := ChallengeService.NewChallengeServiceClient(Clients.ConnChallenge)

	UserController := controller.NewUserController(UserClient, ProblemClient)
	CompilerController := controller.NewCompilerController(NatsClient)
	ProblemController := controller.NewProblemController(ProblemClient, UserClient)
	ChallengeController := controller.NewChallengeController(ChallengeClient, ProblemClient)

	ApiV1 := Router.Group("/api/v1")

	//TEMP: health check and test-inmemory cache
	SetUpTestRoutes(Router)

	SetUpPublicAuthRoutes(ApiV1, UserController)
	SetUpProtectedUserRoutes(ApiV1, UserController, JWTSecret)
	SetUpAdminRoutes(ApiV1, UserController, JWTSecret)
	SetUpCompilerRoutes(ApiV1, CompilerController)
	SetUpProblemRoutes(ApiV1, ProblemController, JWTSecret, UserController)
	SetUpChallengeRoutes(ApiV1, ChallengeController, UserController, JWTSecret)
}

func SetUpTestRoutes(r *gin.Engine) {
	r.GET("/test-cache", func(c *gin.Context) {
		if _, exists := c.Get("cacheInstance"); !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cache instance not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "cache instance ok"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}

// http://localhost:7000/api/v1/auth
func SetUpPublicAuthRoutes(ApiV1 *gin.RouterGroup, UserController *controller.UserController) {
	// http://localhost:7000/api/v1/auth
	Auth := ApiV1.Group("/auth")
	{
		// http://localhost:7000/api/v1/auth/register
		Auth.POST("/register", UserController.RegisterUserHandler)
		// http://localhost:7000/api/v1/auth/login
		Auth.POST("/login", UserController.LoginUserHandler)
		// http://localhost:7000/api/v1/auth/google/login
		Auth.GET("/google/login", UserController.GoogleLoginInitiate)
		// http://localhost:7000/api/v1/auth/google/callback
		Auth.GET("/google/callback", UserController.GoogleLoginCallback)
		// http://localhost:7000/api/v1/auth/token/refresh
		Auth.POST("/token/refresh", UserController.TokenRefreshHandler)
		// http://localhost:7000/api/v1/auth/verify
		Auth.GET("/verify", UserController.VerifyUserHandlerAgainstEmail) //otp verification
		// http://localhost:7000/api/v1/auth/verify/resend
		Auth.GET("/verify/resend", UserController.ResendEmailVerificationHandler)
		// http://localhost:7000/api/v1/auth/password/forgot
		Auth.POST("/password/forgot", UserController.ForgotPasswordHandler)
		// http://localhost:7000/api/v1/auth/password/reset
		Auth.POST("/password/reset", UserController.FinishForgotPasswordHandler)
		// http://localhost:7000/api/v1/auth/2fa/status
		Auth.GET("/2fa/status", UserController.GetTwoFactorAuthStatusHandler)
	}
}

// http://localhost:7000/api/v1/users
func SetUpProtectedUserRoutes(ApiV1 *gin.RouterGroup, UserController *controller.UserController, JWTSecret string) {
	// http://localhost:7000/api/v1/users
	Users := ApiV1.Group("/users")
	// http://localhost:7000/api/v1/users
	UsersPublic := Users.Group("")
	{
		// http://localhost:7000/api/v1/users/public/profile
		UsersPublic.GET("/public/profile", UserController.GetUserProfilePublicHandler)
		// http://localhost:7000/api/v1/users/username/available
		UsersPublic.GET("/username/available", UserController.UserAvailable)

		//http://localhost:7000/api/v1/users/metadata/bulk
		UsersPublic.GET("/metadata/bulk", UserController.GetUsersMetadataBulkList)
	}
	// http://localhost:7000/api/v1/users (protected)
	UsersPrivate := Users.Group("")
	UsersPrivate.Use(
		middleware.JWTAuthMiddleware(JWTSecret),
		middleware.RoleAuthMiddleware(middleware.RoleUser, middleware.RoleAdmin),
		middleware.UserBanCheckMiddleware(UserController.GetUserClient()),
	)

	UsersPrivate.GET("/check-token", UserController.CheckToken)
	{
		// http://localhost:7000/api/v1/users/profile
		Profile := UsersPrivate.Group("/profile")
		{
			// http://localhost:7000/api/v1/users/profile
			Profile.GET("", UserController.GetUserProfileHandler)
			// http://localhost:7000/api/v1/users/profile/update
			Profile.PUT("/update", UserController.UpdateProfileHandler)
			// http://localhost:7000/api/v1/users/profile/image
			Profile.PATCH("/image", UserController.UpdateProfileImageHandler)
			// http://localhost:7000/api/v1/users/profile/ban-history
			Profile.GET("/ban-history", UserController.BanHistoryHandler)
		}
		// http://localhost:7000/api/v1/users/follow
		Follow := UsersPrivate.Group("/follow")
		{
			// http://localhost:7000/api/v1/users/follow
			Follow.POST("", UserController.FollowUserHandler)
			// http://localhost:7000/api/v1/users/follow
			Follow.DELETE("", UserController.UnfollowUserHandler)
			// http://localhost:7000/api/v1/users/follow/following
			Follow.GET("/following", UserController.GetFollowingHandler)
			// http://localhost:7000/api/v1/users/follow/followers
			Follow.GET("/followers", UserController.GetFollowersHandler)
			// http://localhost:7000/api/v1/users/follow/check
			Follow.GET("/check", UserController.GetFollowFollowingCheckHandler)
		}
		// http://localhost:7000/api/v1/users/security
		Security := UsersPrivate.Group("/security")
		{
			// http://localhost:7000/api/v1/users/security/password/change
			Security.POST("/password/change", UserController.ChangePasswordHandler)
			// http://localhost:7000/api/v1/users/security/2fa/setup
			Security.POST("/2fa/setup", UserController.SetUpTwoFactorAuthHandler)
			// http://localhost:7000/api/v1/users/security/2fa/verify
			Security.POST("/2fa/verify", UserController.VerifyTwoFactorAuth)
			// http://localhost:7000/api/v1/users/security/2fa/setup
			Security.DELETE("/2fa/setup", UserController.DisableTwoFactorAuthHandler)
		}
		// http://localhost:7000/api/v1/users/search
		UsersPrivate.GET("/search", UserController.SearchUsersHandler)
		// http://localhost:7000/api/v1/users/logout
		UsersPrivate.POST("/logout", UserController.LogoutUserHandler)
	}
}

// http://localhost:7000/api/v1/admin
func SetUpAdminRoutes(ApiV1 *gin.RouterGroup, UserController *controller.UserController, JWTSecret string) {
	// http://localhost:7000/api/v1/admin
	AdminRoot := ApiV1.Group("/admin")
	{
		// http://localhost:7000/api/v1/admin
		AdminPublic := AdminRoot.Group("")
		{
			// http://localhost:7000/api/v1/admin/login
			AdminPublic.POST("/login", UserController.LoginAdminHandler)
		}
		// http://localhost:7000/api/v1/admin/users
		AdminUsers := AdminRoot.Group("/users")
		AdminUsers.Use(
			middleware.JWTAuthMiddleware(JWTSecret),
			middleware.RoleAuthMiddleware(middleware.RoleAdmin),
		)
		{
			// http://localhost:7000/api/v1/admin/users
			AdminUsers.GET("", UserController.GetAllUsersHandler)
			// http://localhost:7000/api/v1/admin/users
			AdminUsers.POST("", UserController.CreateUserAdminHandler)
			// http://localhost:7000/api/v1/admin/users/update
			AdminUsers.PUT("/update", UserController.UpdateUserAdminHandler)
			// http://localhost:7000/api/v1/admin/users/soft-delete
			AdminUsers.DELETE("/soft-delete", UserController.SoftDeleteUserAdminHandler)
			// http://localhost:7000/api/v1/admin/users/verify
			AdminUsers.POST("/verify", UserController.VerifyAdminUserHandler)
			// http://localhost:7000/api/v1/admin/users/unverify
			AdminUsers.POST("/unverify", UserController.UnverifyUserHandler)
			// http://localhost:7000/api/v1/admin/users/ban
			AdminUsers.POST("/ban", UserController.BanUserHandler)
			// http://localhost:7000/api/v1/admin/users/unban
			AdminUsers.POST("/unban", UserController.UnbanUserHandler)
			// http://localhost:7000/api/v1/admin/users/ban-history
			AdminUsers.GET("/ban-history", UserController.BanHistoryHandler)
		}
	}
}

// http://localhost:7000/api/v1/compile
func SetUpCompilerRoutes(ApiV1 *gin.RouterGroup, CompilerController *controller.CompilerController) {
	// http://localhost:7000/api/v1/compile
	Compiler := ApiV1.Group("")
	{
		// http://localhost:7000/api/v1/compile
		Compiler.POST("/compile", CompilerController.CompileCodeHandler)
	}
}

// http://localhost:7000/api/v1/problems
func SetUpProblemRoutes(ApiV1 *gin.RouterGroup, ProblemController *controller.ProblemController, JWTSecret string, UserController *controller.UserController) {
	// http://localhost:7000/api/v1/problems
	Problems := ApiV1.Group("/problems")
	// http://localhost:7000/api/v1/problems
	ProblemsPublic := Problems.Group("")
	{
		// http://localhost:7000/api/v1/problems/list
		ProblemsPublic.GET("/list", ProblemController.ListProblemsHandler)
		// http://localhost:7000/api/v1/problems/metadata
		ProblemsPublic.GET("/metadata", ProblemController.GetProblemByIDSlugHandler)
		// http://localhost:7000/api/v1/problems/metadata/list
		ProblemsPublic.GET("/metadata/list", ProblemController.GetProblemMetadataListHandler)
		// http://localhost:7000/api/v1/problems/leaderboard/top10
		ProblemsPublic.GET("/leaderboard/top10", ProblemController.GetTopKGlobalController)
		// http://localhost:7000/api/v1/problems/leaderboard/top10/entity
		ProblemsPublic.GET("/leaderboard/top10/entity", ProblemController.GetTopKEntityController)
		// http://localhost:7000/api/v1/problems/languages
		ProblemsPublic.GET("/languages", ProblemController.GetLanguageSupportsHandler)

		// http://localhost:7000/api/v1/problems/bulk/metadata - queryarray [problemIds]
		ProblemsPublic.GET("/bulk/metadata", ProblemController.GetBulkProblemMetadata)

		// http://localhost:7000/api/v1/problems/execute
		ProblemsPublic.POST("/execute", ProblemController.RunUserCodeProblemHandler)
		// http://localhost:7000/api/v1/problems/submission/history
		ProblemsPublic.GET("/submission/history", ProblemController.GetSubmissionHistoryOptionalProblemId)
		// http://localhost:7000/api/v1/problems/stats
		ProblemsPublic.GET("/stats", ProblemController.GetProblemStatistics)
		// http://localhost:7000/api/v1/problems/activity
		ProblemsPublic.GET("/activity", ProblemController.GetMonthlyActivityHeatmapController)
		// http://localhost:7000/api/v1/problems/leaderboard/data
		ProblemsPublic.GET("/leaderboard/data", ProblemController.GetLeaderboardDataController)
		// http://localhost:7000/api/v1/problems/user/done
		ProblemsPublic.GET("/user/done", ProblemController.ProblemIDsDoneByUserID)
		// http://localhost:7000/api/v1/problems/verify/bulk

		ProblemsPublic.POST("/verify/bulk", ProblemController.VerifyProblemExistenceBulk)
		// http://localhost:7000/api/v1/problems/random
		ProblemsPublic.GET("/random", ProblemController.RandomProblemIDsGenWithDifficultyRatio)

		//http://localhost:7000/api/v1/problems/count
		ProblemsPublic.GET("/count", ProblemController.ProblemCountMetadata)
	}
	// http://localhost:7000/api/v1/problems (protected)
	ProblemsPrivate := Problems.Group("")
	ProblemsPrivate.Use(
		middleware.JWTAuthMiddleware(JWTSecret),
		middleware.RoleAuthMiddleware(middleware.RoleAdmin),
	)
	{
		// http://localhost:7000/api/v1/problems/list/all
		ProblemsPrivate.GET("/list/all", ProblemController.ListProblemsHandler)
		// http://localhost:7000/api/v1/problems/
		ProblemsPrivate.POST("/", ProblemController.CreateProblemHandler)
		// http://localhost:7000/api/v1/problems/
		ProblemsPrivate.PUT("/", ProblemController.UpdateProblemHandler)
		// http://localhost:7000/api/v1/problems/
		ProblemsPrivate.DELETE("/", ProblemController.DeleteProblemHandler)
		// http://localhost:7000/api/v1/problems/
		ProblemsPrivate.GET("/", ProblemController.GetProblemHandler)
		// http://localhost:7000/api/v1/problems/testcases
		ProblemsPrivate.POST("/testcases", ProblemController.AddTestCasesHandler)
		// http://localhost:7000/api/v1/problems/testcases/single
		ProblemsPrivate.DELETE("/testcases/single", ProblemController.DeleteTestCaseHandler)
		// http://localhost:7000/api/v1/problems/language
		ProblemsPrivate.POST("/language", ProblemController.AddLanguageSupportHandler)
		// http://localhost:7000/api/v1/problems/language
		ProblemsPrivate.PUT("/language", ProblemController.UpdateLanguageSupportHandler)
		// http://localhost:7000/api/v1/problems/language
		ProblemsPrivate.DELETE("/language", ProblemController.RemoveLanguageSupportHandler)
		// http://localhost:7000/api/v1/problems/validate
		ProblemsPrivate.GET("/validate", ProblemController.FullValidationByProblemIDHandler)
		// http://localhost:7000/api/v1/problems/leaderboard/rank
		ProblemsPrivate.GET("/leaderboard/rank", ProblemController.GetUserRankController)
	}

}

// TODO: migrate to challenge service
func SetUpChallengeRoutes(apiV1 *gin.RouterGroup, challengeController *controller.ChallengeController, userController *controller.UserController, jwtSecret string) {
	// route group for /api/v1/challenges
	challenges := apiV1.Group("/challenges")

	// protected routes
	challengesPrivate := challenges.Group("")
	challengesPrivate.Use(
		middleware.JWTAuthMiddleware(jwtSecret),
		middleware.RoleAuthMiddleware(middleware.RoleUser, middleware.RoleAdmin),
		middleware.UserBanCheckMiddleware(userController.GetUserClient()),
	)
	{
		// POST /api/v1/challenges/
		challengesPrivate.POST("", challengeController.CreateChallenge)
		// POST /api/v1/challenges/abandon
		challengesPrivate.POST("/abandon", challengeController.AbandonChallenge)
		// GET /api/v1/challenges/history
		challengesPrivate.GET("/history", challengeController.GetChallengeHistory)

		//GetOwnersActiveChallenges
		challengesPrivate.GET("/owner/open", challengeController.GetOwnersActiveChallenges)
		// GET /api/v1/challenges/:challengeId
		challengesPrivate.GET("/:challengeId", challengeController.GetChallengeByID)
	}

	// GET /api/v1/challenges/public/open?page=1&page_size=10&is_private=false
	challenges.GET("/public/open", challengeController.GetActiveOpenChallenges) //page page_size
}
