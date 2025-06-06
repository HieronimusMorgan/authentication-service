package controller

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/services"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(c *gin.Context)
	RegisterDeviceToken(c *gin.Context)
	Login(c *gin.Context)
	LoginPhoneNumber(c *gin.Context)
	ResetPassword(c *gin.Context)
	ChangeDeviceID(c *gin.Context)
	VerifyDeviceID(c *gin.Context)
	VerifyPinCode(c *gin.Context)
	ChangePinCode(c *gin.Context)
	ForgetPinCode(c *gin.Context)
	RefreshToken(c *gin.Context)
	RegisterInternalToken(c *gin.Context)
	UpdateRole(ctx *gin.Context)
	GetListUser(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	GenerateCredentialKey(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type authController struct {
	AuthService services.AuthService
	UserSession services.UsersSessionService
	JWTService  utils.JWTService
}

func NewAuthController(serviceAuth services.AuthService, serviceSession services.UsersSessionService, jwtService utils.JWTService) AuthController {
	return authController{AuthService: serviceAuth, UserSession: serviceSession, JWTService: jwtService}
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
func (h authController) Register(c *gin.Context) {
	var req in.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}
	var deviceID = c.GetHeader("Device-Type")

	if deviceID != "WEB" && deviceID != "MOBILE" {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid or missing Device-Type", nil)
		return
	}

	user, err := h.AuthService.Register(&req, deviceID)
	if err != nil {
		handleErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	errSession := h.UserSession.AddUserSession(user.UserID, user.Token,
		user.RefreshToken, c.ClientIP(), "WEB")

	if errSession != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to create user session", err)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "User registered successfully", user)
}

func (h authController) RegisterDeviceToken(c *gin.Context) {
	var req struct {
		DeviceToken string `json:"device_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	errs := h.AuthService.RegisterDeviceToken(&req, token.ClientID)
	if errs != nil {
		handleErrorResponse(c, http.StatusBadRequest, errs.Error(), nil)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Device token registered successfully", nil)
}

func (h authController) Login(c *gin.Context) {
	var req in.LoginRequest
	var deviceID = c.GetHeader("Device-Type")

	if deviceID != "WEB" && deviceID != "MOBILE" {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid or missing Device-Type", nil)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user, err := h.AuthService.Login(&req, deviceID)
	if err != nil {
		handleErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	errSession := h.UserSession.AddUserSession(user.(out.LoginResponse).UserID, user.(out.LoginResponse).Token,
		user.(out.LoginResponse).RefreshToken, c.ClientIP(), deviceID)

	if errSession != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to create user session", err)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Login successful", user)
}

func (h authController) LoginPhoneNumber(c *gin.Context) {
	var req in.LoginPhoneNumber
	var deviceID = c.GetHeader("Device-Type")

	if deviceID != "WEB" && deviceID != "MOBILE" {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid or missing Device-Type", nil)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user, errs := h.AuthService.LoginPhoneNumber(&req, deviceID)
	if errs != nil {
		handleErrorResponse(c, http.StatusBadRequest, errs.Error(), nil)
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

func (h authController) ResetPassword(c *gin.Context) {
	var req struct {
		NewPassword     string `json:"new_password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
		RequestID       string `json:"request_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	errs := h.AuthService.ResetPassword(&req)
	if errs != nil {
		handleErrorResponse(c, http.StatusBadRequest, errs.Error(), nil)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}

func (h authController) ChangeDeviceID(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		DeviceID    string `json:"device_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	data, errs := h.AuthService.ChangeDeviceID(&req)
	if errs != nil {
		handleErrorResponse(c, http.StatusBadRequest, errs.Error(), nil)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Device ID changed successfully", data)
}

func (h authController) VerifyDeviceID(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id" binding:"required"`
		PinCode   string `json:"pin_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	data, errs := h.AuthService.VerifyDeviceID(&req)
	if errs != nil {
		handleErrorResponse(c, http.StatusBadRequest, errs.Error(), nil)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Device ID verified successfully", data)
}

func (h authController) VerifyPinCode(c *gin.Context) {
	var req struct {
		PinCode string `json:"pin_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	clientID, err := h.AuthService.VerifyPinCode(&req, token.ClientID)
	if err != nil {
		handleErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Pin verified successfully", clientID)
}

func (h authController) ChangePinCode(c *gin.Context) {
	var req struct {
		OldPinCode string `json:"old_pin_code" binding:"required"`
		NewPinCode string `json:"new_pin_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	errs := h.AuthService.ChangePinCode(&req, token.ClientID)
	if errs != nil {
		handleErrorResponse(c, http.StatusBadRequest, errs.Error(), nil)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Pin changed successfully", nil)
}

func (h authController) ForgetPinCode(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required"`
		PinCode string `json:"pin_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	errs := h.AuthService.ForgetPinCode(&req, token.ClientID)
	if errs != nil {
		handleErrorResponse(c, http.StatusBadRequest, errs.Error(), nil)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Pin reset successfully", nil)
}

func (h authController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	newToken, errs := h.AuthService.RefreshToken(&req, token.ClientID)
	if errs != nil {
		handleErrorResponse(c, http.StatusBadRequest, errs.Error(), nil)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Token refreshed successfully", newToken)
}

func (h authController) RegisterInternalToken(c *gin.Context) {
	var req struct {
		ResourceName string `json:"resource_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, errs := h.AuthService.RegisterInternalToken(&req)
	if errs != nil {
		response.SendResponse(c, http.StatusBadRequest, errs.Error(), nil, errs)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Internal token registered successfully", token)
}

func (h authController) UpdateRole(ctx *gin.Context) {
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

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	errs := h.AuthService.UpdateRole(userID, req.RoleID, token.ClientID)
	if errs != nil {
		response.SendResponse(ctx, http.StatusBadRequest, errs.Error(), nil, errs)
		return
	}

	response.SendResponse(ctx, 200, "Role updated successfully", nil, nil)
}

func (h authController) GetListUser(ctx *gin.Context) {
	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	users, errs := h.AuthService.GetListUser(token.ClientID)
	if errs != nil {
		response.SendResponse(ctx, http.StatusBadRequest, errs.Error(), nil, errs)
		return
	}

	response.SendResponse(ctx, 200, "List user retrieved successfully", users, nil)
}

func (h authController) GetUserByID(ctx *gin.Context) {
	userID, err := utils.ConvertToUint(ctx.Param("id"))
	if err != nil {
		response.SendResponse(ctx, 400, "User ID must be a number", nil, err)
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	user, errs := h.AuthService.GetUserByID(userID, token.ClientID)
	if errs != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, errs)
		return
	}

	response.SendResponse(ctx, 200, "User retrieved successfully", user, nil)
}

func (h authController) ForgotPassword(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err)
		return
	}

	errs := h.AuthService.RequestForgotPassword(&req)
	if errs != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, errs)
		return
	}

	response.SendResponse(ctx, 200, "Password reset request successful", nil, nil)
}

func (h authController) ChangePassword(ctx *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendResponse(ctx, 400, "Invalid request", nil, err)
		return
	}

	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	errs := h.AuthService.ChangePassword(&req, token.ClientID)
	if errs != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, errs)
		return
	}

	response.SendResponse(ctx, 200, "Password changed successfully", nil, nil)
}

func (h authController) GenerateCredentialKey(ctx *gin.Context) {
	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}
	credentialKey, err := h.AuthService.GenerateCredentialKey(token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to get credential key", nil, err.Error())
		return
	}

	response.SendResponse(ctx, 200, "Get credential key successfully", credentialKey, nil)
}
func (h authController) Logout(ctx *gin.Context) {
	token, exist := utils.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := h.UserSession.LogoutSession(token.UserID)
	if err != nil {
		response.SendResponse(ctx, 400, "User Session", nil, err.Error())
		return
	}

	response.SendResponse(ctx, 200, "Logout successful", nil, nil)
}
