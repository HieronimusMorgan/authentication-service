package handler

import (
	"authentication/internal/dto/in"
	"authentication/internal/services"
	"authentication/package/response"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{AuthService: services.NewAuthService(db)}
}

// Helper for centralized error response
func handleErrorResponse(c *gin.Context, status int, message string, err error) {
	if err != nil {
		response.SendResponse(c, status, message, nil, err.Error())
		return
	}
}

// Helper for centralized success response
func handleSuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	response.SendResponse(c, status, message, data, nil)
}

// Register a new user
func (h AuthHandler) Register(c *gin.Context) {
	var req in.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user, err := h.AuthService.Register(&req)
	if err != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	handleSuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

func (h AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user, err := h.AuthService.Login(&req)
	if err != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Login successful", user)
}

func (h AuthHandler) GetProfile(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		handleErrorResponse(c, http.StatusUnauthorized, "Authorization token is required", nil)
		return
	}

	user, err := h.AuthService.GetProfile(token)
	if err != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to get profile", err)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}

func (h AuthHandler) RegisterInternalToken(c *gin.Context) {
	var req struct {
		ResourceName string `json:"resource_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, err := h.AuthService.RegisterInternalToken(&req)
	if err != nil {
		handleErrorResponse(c, http.StatusInternalServerError, "Failed to register internal token", err)
		return
	}

	handleSuccessResponse(c, http.StatusCreated, "Internal token registered successfully", token)
}
