package routes

import (
	"authentication/config"
	"authentication/internal/controller"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, middleware config.Middleware, authController controller.AuthController) {
	public := r.Group("/v1")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
		public.POST("/login-phone", authController.LoginPhoneNumber)
	}

	protected := r.Group("/v1")
	protected.Use(middleware.AuthMiddleware.Handler())
	{
		protected.POST("/verify-pin", authController.VerifyPinCode)
		protected.POST("/change-password", authController.ChangePassword)
		protected.POST("/change-pin", authController.ChangePinCode)
		protected.POST("/forget-pin", authController.ForgetPinCode)
		protected.GET("/logout", authController.Logout)
	}

	admin := r.Group("/v1")
	admin.Use(middleware.AdminMiddleware.Handler())
	{
		admin.GET("/users", authController.GetListUser)
		admin.POST("/user/update-role/:id", authController.UpdateRole)
	}
}
