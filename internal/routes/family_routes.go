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
		protected.POST("/add-permission", familyController.AddFamilyMemberPermission)
		protected.POST("/remove-permission", familyController.RemoveFamilyMemberPermission)
		protected.GET("/member/:id", familyController.GetFamilyMembers)
		protected.POST("/delete/:id", familyController.RemoveMemberFamily)
	}

}
