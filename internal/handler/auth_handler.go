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
	s := services.NewAuthService(db)
	return &AuthHandler{AuthService: s}
}

func (h AuthHandler) Register(c *gin.Context) {
	var req in.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	user, err := h.AuthService.Register(&req)
	if err != nil {
		response.SendResponse(c, http.StatusInternalServerError, "Failed to register user", nil, err)
		return
	}

	response.SendResponse(c, http.StatusCreated, "User registered successfully", user, nil)
}

func (h AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	user, err := h.AuthService.Login(&req)
	if err != nil {
		response.SendResponse(c, http.StatusInternalServerError, "Failed to login", nil, err)
		return
	}

	response.SendResponse(c, http.StatusOK, "Login success", user, nil)
}

func (h AuthHandler) GetProfile(c *gin.Context) {
	token := c.GetHeader("Authorization")

	user, err := h.AuthService.GetProfile(token)
	if err != nil {
		response.SendResponse(c, http.StatusInternalServerError, "Failed to get profile", nil, err)
		return
	}

	response.SendResponse(c, http.StatusOK, "Get profile success", user, nil)
}

func (h AuthHandler) RegisterInternalToken(context *gin.Context) {
	var req struct {
		ResourceName string `json:"resource_name" binding:"required"`
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	token, err := h.AuthService.RegisterInternalToken(&req)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Failed to register user", nil, err)
		return
	}
	response.SendResponse(context, http.StatusCreated, "User registered successfully", token, nil)
}
