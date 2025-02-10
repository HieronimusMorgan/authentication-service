package handler

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/services"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	RegisterInternalToken(c *gin.Context)
	DeleteUser(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	GetListUser(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type authHandler struct {
	AuthService services.AuthService
	UserSession services.UsersSessionService
	JWTService  utils.JWTService
}

func NewAuthHandler(serviceAuth services.AuthService, serviceSession services.UsersSessionService, jwtService utils.JWTService) AuthHandler {
	return authHandler{AuthService: serviceAuth, UserSession: serviceSession, JWTService: jwtService}
}

// Helper for centralized error response
func handleErrorResponse(c *gin.Context, status int, message string, err interface{}) {
	response.SendResponse(c, status, message, nil, err)
	return
}

// Helper for centralized success response
func handleSuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	response.SendResponse(c, status, message, data, nil)
}

// Register a new user
func (h authHandler) Register(c *gin.Context) {
	var errs = response.ErrorResponse{}
	var req in.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user, errs := h.AuthService.Register(&req)
	if errs.Message != "" {
		handleErrorResponse(c, errs.Code, errs.Message, nil)
		return
	}

	err := h.UserSession.AddUserSession(user.UserID, user.Token,
		user.RefreshToken, c.ClientIP(), "WEB")

	if err != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to create user session", err)
		return
	}

	handleSuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

func (h authHandler) Login(c *gin.Context) {
	var errs = response.ErrorResponse{}
	var req in.LoginRequest
	var deviceID = c.GetHeader("Device-ID")

	if deviceID != "WEB" && deviceID != "MOBILE" {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid or missing Device-ID", nil)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user, errs := h.AuthService.Login(&req)
	if errs.Message != "" {
		handleErrorResponse(c, errs.Code, errs.Message, nil)
		return
	}

	err := h.UserSession.AddUserSession(user.(out.LoginResponse).UserID, user.(out.LoginResponse).Token,
		user.(out.LoginResponse).RefreshToken, c.ClientIP(), deviceID)

	if err != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to create user session", err)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Login successful", user)
}

func (h authHandler) GetProfile(c *gin.Context) {
	var errs = response.ErrorResponse{}
	token := c.GetHeader("Authorization")
	if token == "" {
		handleErrorResponse(c, http.StatusUnauthorized, "Authorization token is required", nil)
		return
	}

	user, errs := h.AuthService.GetProfile(token)
	if errs.Message != "" {
		response.SendResponse(c, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}

func (h authHandler) RegisterInternalToken(c *gin.Context) {
	var errs = response.ErrorResponse{}
	var req struct {
		ResourceName string `json:"resource_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, errs := h.AuthService.RegisterInternalToken(&req)
	if errs.Message != "" {
		response.SendResponse(c, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	handleSuccessResponse(c, http.StatusCreated, "Internal token registered successfully", token)
}

func (h authHandler) DeleteUser(ctx *gin.Context) {
	var errs = response.ErrorResponse{}
	userID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Resource ID must be a number", nil, err)
		return
	}

	token, err := extractClaims(ctx, h.JWTService)
	if err != nil {
		return
	}

	errs = h.AuthService.DeleteUserById(userID, token.ClientID)
	if errs.Message != "" {
		response.SendResponse(ctx, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	response.SendResponse(ctx, 200, "User deleted successfully", nil, nil)
}

func (h authHandler) UpdateRole(ctx *gin.Context) {
	var errs = response.ErrorResponse{}
	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err)
		return
	}

	userID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "User ID must be a number", nil, err)
		return
	}

	token, err := extractClaims(ctx, h.JWTService)
	if err != nil {
		return
	}

	errs = h.AuthService.UpdateRole(userID, req.RoleID, token.ClientID)
	if errs.Message != "" {
		response.SendResponse(ctx, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	response.SendResponse(ctx, 200, "Role updated successfully", nil, nil)
}

func (h authHandler) GetListUser(ctx *gin.Context) {
	var errs = response.ErrorResponse{}
	token, err := extractClaims(ctx, h.JWTService)
	if err != nil {
		return
	}

	users, errs := h.AuthService.GetListUser(token.ClientID)
	if errs.Message != "" {
		response.SendResponse(ctx, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	response.SendResponse(ctx, 200, "List user retrieved successfully", users, nil)
}

func (h authHandler) ChangePassword(ctx *gin.Context) {
	var errs = response.ErrorResponse{}
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err)
		return
	}

	token, err := extractClaims(ctx, h.JWTService)
	if err != nil {
		return
	}

	errs = h.AuthService.ChangePassword(&req, token.ClientID)
	if errs != (response.ErrorResponse{}) {
		response.SendResponse(ctx, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	response.SendResponse(ctx, 200, "Password changed successfully", nil, nil)
}

func (h authHandler) Logout(ctx *gin.Context) {
	token, err := extractClaims(ctx, h.JWTService)
	if err != nil {
		return
	}

	err = h.UserSession.LogoutSession(token.UserID)
	if err != nil {
		response.SendResponse(ctx, 400, "User Session", nil, err.Error())
		return
	}

	response.SendResponse(ctx, 200, "Logout successful", nil, nil)
}
