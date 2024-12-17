package handler

import (
	"authentication/internal/services"
	"authentication/internal/utils"
	"authentication/package/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResourceHandler struct {
	ResourceService *services.ResourceService
}

func NewResourceHandler(db *gorm.DB) *ResourceHandler {
	return &ResourceHandler{ResourceService: services.NewResourceService(db)}
}

func extractClaims(context *gin.Context) (utils.TokenClaims, error) {
	token, err := utils.ExtractClaims(context.GetHeader("Authorization"))
	if err != nil {
		response.SendResponse(context, 401, "Unauthorized", nil, err.Error())
	}
	return *token, err
}

func (h ResourceHandler) AddResource(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	token, err := extractClaims(ctx)
	if err != nil {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err.Error())
		return
	}

	resource, err := h.ResourceService.AddResource(&req.Name, &req.Description, token.ClientID)
	response.SendResponse(ctx, 200, "Resource added successfully", resource, err.Error())
}

func (h ResourceHandler) UpdateResource(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	resourceID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, err := extractClaims(ctx)
	if err != nil {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err.Error())
		return
	}

	resource, err := h.ResourceService.UpdateResource(resourceID, &req.Name, &req.Description, token.ClientID)
	response.SendResponse(ctx, 200, "Resource updated successfully", resource, err.Error())
}

func (h ResourceHandler) GetResources(ctx *gin.Context) {
	token, err := extractClaims(ctx)
	if err != nil {
		return
	}

	resources, err := h.ResourceService.GetResources(token.ClientID)
	response.SendResponse(ctx, 200, "Resources retrieved successfully", resources, err.Error())
}

func (h ResourceHandler) AssignResourceToRole(ctx *gin.Context) {
	var req struct {
		RoleID     uint `json:"role_id" binding:"required"`
		ResourceID uint `json:"resource_id" binding:"required"`
	}

	token, err := extractClaims(ctx)
	if err != nil {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err.Error())
		return
	}

	roleResource, err := h.ResourceService.AssignResourceToRole(req.RoleID, req.ResourceID, token.ClientID)
	response.SendResponse(ctx, 200, "Resource assigned to role successfully", roleResource, err.Error())
}

func (h ResourceHandler) GetResourcesById(ctx *gin.Context) {
	resourceID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, err := extractClaims(ctx)
	if err != nil {
		return
	}

	resource, err := h.ResourceService.GetResourceById(resourceID, token.ClientID)
	response.SendResponse(ctx, 200, "Resource retrieved successfully", resource, err.Error())
}

func (h ResourceHandler) DeleteResourceById(ctx *gin.Context) {
	resourceID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, err := extractClaims(ctx)
	if err != nil {
		return
	}

	err = h.ResourceService.DeleteResourceById(resourceID, token.ClientID)
	response.SendResponse(ctx, 200, "Resource deleted successfully", nil, err.Error())
}

func (h ResourceHandler) GetResourceUserById(ctx *gin.Context) {
	resourceID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, err := extractClaims(ctx)
	if err != nil {
		return
	}

	users, err := h.ResourceService.GetResourceUserById(resourceID, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to get users", nil, err.Error())
		return
	}

	response.SendResponse(ctx, 200, "Users retrieved successfully", users, nil)
}
