package controller

import (
	"authentication/internal/dto/in"
	"authentication/internal/services"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetProfile(c *gin.Context)
	UpdateNameUserProfile(c *gin.Context)
	UpdatePhotoUserProfile(c *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userController struct {
	UserService services.UserService
	JWTService  utils.JWTService
}

func NewUserController(serviceUser services.UserService, jwtService utils.JWTService) UserController {
	return userController{UserService: serviceUser, JWTService: jwtService}
}

func (h userController) GetProfile(c *gin.Context) {
	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	user, errs := h.UserService.GetProfile(token.ClientID)
	if errs.Message != "" {
		response.SendResponse(c, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}

func (h userController) UpdateNameUserProfile(c *gin.Context) {
	var req in.UpdateNameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	user, err := h.UserService.UpdateNameUserProfile(&req, token.ClientID)
	if err != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Profile updated successfully", user)
}

func (h userController) UpdatePhotoUserProfile(c *gin.Context) {
	var req in.UpdatePhotoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	user, err := h.UserService.UpdatePhotoUserProfile(&req, token.ClientID)
	if err != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Profile updated successfully", user)
}

func (h userController) DeleteUser(ctx *gin.Context) {
	userID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err)
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	errs := h.UserService.DeleteUserById(userID, token.ClientID)
	if errs.Message != "" {
		response.SendResponse(ctx, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	response.SendResponse(ctx, 200, "User deleted successfully", nil, nil)
}
