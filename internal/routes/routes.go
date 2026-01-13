package routes

import (
	"authentication/backend/internal/config"
	"authentication/backend/internal/handlers"
	"authentication/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	r *gin.Engine,
	authH *handlers.AuthHandler,
	jobH *handlers.Handler,
) {

	// PUBLIC ROUTES
	public := r.Group("/")
	{
		public.POST("/register", authH.Register)
		public.POST("/login", authH.Login)
		public.POST("/forgot-password", authH.ForgotPassword)
		public.POST("/reset-password/:token", authH.ResetPassword)
	}

	// AUTH TEST
	authTest := r.Group("/test/auth")
	authTest.Use(middleware.Auth())
	{
		authTest.GET("", handlers.TestAuth)
	}

	// AUTHENTICATED USER ROUTES
	protected := r.Group("/auth")
	protected.Use(
		middleware.Auth(),
	)
	{
		protected.POST("/convert", jobH.Convert)
		protected.GET("/jobs/:id", jobH.Status)
		protected.GET("/download/:id", jobH.Download)

		protected.POST("/change-password", authH.ChangePassword)
		protected.GET("/test", handlers.TestProtected)
	}

	// ADMIN ROUTES
	admin := r.Group("/admin")
	admin.Use(
		middleware.Auth(),
		middleware.AdminOnly(),
		middleware.RateLimit(config.Redis),
	)
	{
		admin.GET("/blocked-ips", handlers.ListBlockedIPs(config.Redis))
		admin.DELETE("/blocked-ips/:ip", handlers.UnblockIP(config.Redis))
	}
}
