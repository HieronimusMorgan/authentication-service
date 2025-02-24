package controller

import (
	"authentication/internal/services"
	"authentication/internal/utils"
	"authentication/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResourceController interface {
	AddResource(ctx *gin.Context)
	UpdateResource(ctx *gin.Context)
	GetResources(ctx *gin.Context)
	AssignResourceToRole(ctx *gin.Context)
	GetResourcesById(ctx *gin.Context)
	DeleteResourceById(ctx *gin.Context)
	GetResourceUserById(ctx *gin.Context)
	GetResourceRoles(ctx *gin.Context)
}

type resourceController struct {
	ResourceService services.ResourceService
	JWTService      utils.JWTService
}

func NewResourceController(resourceService services.ResourceService, jwtService utils.JWTService) ResourceController {
	return resourceController{ResourceService: resourceService, JWTService: jwtService}
}

func (h resourceController) AddResource(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err.Error())
		return
	}

	resource, err := h.ResourceService.AddResource(&req.Name, &req.Description, token.ClientID)
	response.SendResponse(ctx, 200, "Resource added successfully", resource, err.Error())
}

func (h resourceController) UpdateResource(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	resourceID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err.Error())
		return
	}

	resource, err := h.ResourceService.UpdateResource(resourceID, &req.Name, &req.Description, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 400, "Failed to update resource", nil, err.Error())
		return
	}
	response.SendResponse(ctx, 200, "Resource updated successfully", resource, nil)
}

func (h resourceController) GetResources(ctx *gin.Context) {
	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	resources, err := h.ResourceService.GetResources(token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to get resources", nil, err.Error())
		return
	}
	response.SendResponse(ctx, 200, "Resources retrieved successfully", resources, nil)
}

func (h resourceController) AssignResourceToRole(ctx *gin.Context) {
	var req struct {
		RoleID     uint `json:"role_id" binding:"required"`
		ResourceID uint `json:"resource_id" binding:"required"`
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err.Error())
		return
	}

	roleResource, err := h.ResourceService.AssignResourceToRole(req.RoleID, req.ResourceID, token.ClientID)
	response.SendResponse(ctx, 200, "Resource assigned to role successfully", roleResource, err.Error())
}

func (h resourceController) GetResourcesById(ctx *gin.Context) {
	resourceID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	resource, err := h.ResourceService.GetResourceById(resourceID, token.ClientID)
	response.SendResponse(ctx, 200, "Resource retrieved successfully", resource, err.Error())
}

func (h resourceController) DeleteResourceById(ctx *gin.Context) {
	resourceID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = h.ResourceService.DeleteResourceById(resourceID, token.ClientID)
	response.SendResponse(ctx, 200, "Resource deleted successfully", nil, err.Error())
}

func (h resourceController) GetResourceUserById(ctx *gin.Context) {
	resourceID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	users, err := h.ResourceService.GetResourceUserById(resourceID, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to get users", nil, err.Error())
		return
	}

	response.SendResponse(ctx, 200, "Users retrieved successfully", users, nil)
}

func (h resourceController) GetResourceRoles(context *gin.Context) {
	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	roles, err := h.ResourceService.GetResourceRoles(token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to get roles", nil, err.Error())
		return
	}

	response.SendResponse(context, 200, "Roles retrieved successfully", roles, nil)
}
