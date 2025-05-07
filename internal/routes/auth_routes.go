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
		public.POST("/change-device", authController.ChangeDeviceID)
		public.POST("/verify-device", authController.VerifyDeviceID)
	}

	protected := r.Group("/v1")
	protected.Use(middleware.AuthMiddleware.Handler())
	{
		protected.POST("/register-device-token", authController.RegisterDeviceToken)
		protected.GET("/credential-key", authController.GenerateCredentialKey)
		protected.POST("/verify-pin", authController.VerifyPinCode)
		protected.POST("/change-password", authController.ChangePassword)
		protected.POST("/change-pin", authController.ChangePinCode)
		protected.POST("/forget-pin", authController.ForgetPinCode)
		protected.GET("/logout", authController.Logout)
		protected.POST("/refresh-token", authController.RefreshToken)
	}

	admin := r.Group("/v1")
	admin.Use(middleware.AdminMiddleware.Handler())
	{
		admin.GET("/users", authController.GetListUser)
		admin.GET("/users/:id", authController.GetUserByID)
		admin.POST("/user/update-role/:id", authController.UpdateRole)
	}
}
