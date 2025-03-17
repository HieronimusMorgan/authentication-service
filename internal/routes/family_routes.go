package routes

import (
	"authentication/config"
	"authentication/internal/controller"
	"github.com/gin-gonic/gin"
)

func FamilyRoutes(r *gin.Engine, middleware config.Middleware, familyController controller.FamilyController) {

	protected := r.Group("/v1/family")
	protected.Use(middleware.AuthMiddleware.Handler())
	{
		protected.POST("/create", familyController.CreateFamily)
		protected.POST("/add-member", familyController.AddMemberFamily)
	}

}
