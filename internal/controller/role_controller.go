package controller

import (
	"authentication/internal/services"
	"authentication/internal/utils"
	"authentication/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoleController interface {
	AddRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	GetListRole(ctx *gin.Context)
	GetListUserRole(ctx *gin.Context)
	GetRoleById(ctx *gin.Context)
	DeleteRoleById(ctx *gin.Context)
	GetListRoleUsers(ctx *gin.Context)
}

type roleController struct {
	RoleService services.RoleService
	JWTService  utils.JWTService
}

func NewRoleController(db services.RoleService, jwtService utils.JWTService) RoleController {
	return roleController{RoleService: db, JWTService: jwtService}
}

func extractAndValidateToken(context *gin.Context, service utils.JWTService) (utils.TokenClaims, error) {
	token, err := service.ExtractClaims(context.GetHeader("Authorization"))
	if err != nil {
		response.SendResponse(context, 401, "Unauthorized", nil, err.Error())
	}
	return *token, err
}

func (h roleController) AddRole(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	role, err := h.RoleService.RegisterRole(&req, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to register role", nil, err.Error())
		return
	}
	response.SendResponse(ctx, 201, "Role registered successfully", role, nil)
}

func (h roleController) UpdateRole(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"optional"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err.Error())
		return
	}

	id := ctx.Param("id")
	if id == "" {
		response.SendResponse(ctx, 400, "Role ID is required", nil, nil)
		return
	}

	roleID, err := utils.ConvertToUint(id)
	if err != nil {
		response.SendResponse(ctx, 400, "Role ID must be a number", nil, err.Error())
		return
	}

	if req.Name == "" {
		response.SendResponse(ctx, 400, "Role name is required", nil, nil)
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	role, err := h.RoleService.UpdateRole(roleID, &req, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to update role", nil, err.Error())
		return
	}
	response.SendResponse(ctx, 200, "Role updated successfully", role, nil)
}

func (h roleController) GetListRole(ctx *gin.Context) {
	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	roles, err := h.RoleService.GetListRole(token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to get list of roles", nil, err.Error())
		return
	}

	response.SendResponse(ctx, 200, "Roles retrieved successfully", roles, nil)
}

func (h roleController) GetListUserRole(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.SendResponse(ctx, 400, "Role ID is required", nil, nil)
		return
	}

	roleID, err := utils.ConvertToUint(id)
	if err != nil {
		response.SendResponse(ctx, 400, "Role ID must be a number", nil, err.Error())
		return
	}

	pageIndex, pageSize, err := utils.GetPageIndexPageSize(ctx)
	if err != nil {
		response.SendResponse(ctx, 400, "Invalid page index or page size", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	userRoles, total, err := h.RoleService.GetListUserRole(token.ClientID, roleID, pageIndex, pageSize)
	if err != nil {
		response.SendResponseList(ctx, 500, "Failed to get list of user roles", response.PagedData{
			Total:     total,
			PageIndex: pageIndex,
			PageSize:  pageSize,
			Items:     nil,
		}, err.Error())
		return
	}

	response.SendResponseList(ctx, 200, "User roles retrieved successfully", response.PagedData{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Items:     userRoles,
	}, nil)
}

func (h roleController) GetRoleById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.SendResponse(ctx, 400, "Role ID is required", nil, nil)
		return
	}

	roleID, err := utils.ConvertToUint(id)
	if err != nil {
		response.SendResponse(ctx, 400, "Role ID must be a number", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	role, err := h.RoleService.GetRoleById(roleID, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to retrieve role", nil, err.Error())
		return
	}

	response.SendResponse(ctx, 200, "Role retrieved successfully", role, nil)
}

func (h roleController) DeleteRoleById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.SendResponse(ctx, 400, "Role ID is required", nil, nil)
		return
	}

	roleID, err := utils.ConvertToUint(id)
	if err != nil {
		response.SendResponse(ctx, 400, "Role ID must be a number", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = h.RoleService.DeleteRole(roleID, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to delete role", nil, err.Error())
		return
	}

	response.SendResponse(ctx, 200, "Role deleted successfully", nil, nil)
}

func (h roleController) GetListRoleUsers(ctx *gin.Context) {
	pageIndex, pageSize, err := utils.GetPageIndexPageSize(ctx)
	if err != nil {
		response.SendResponse(ctx, 400, "Invalid page index or page size", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	roleUsers, total, err := h.RoleService.GetListRoleUsers(token.ClientID, pageIndex, pageSize)
	if err != nil {
		response.SendResponseList(ctx, 500, "Failed to get list assets", response.PagedData{
			Total:     total,
			PageIndex: pageIndex,
			PageSize:  pageSize,
			Items:     nil,
		}, err.Error())
		return
	}

	response.SendResponseList(ctx, 200, "Get list assets successfully", response.PagedData{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Items:     roleUsers,
	}, nil)
}
