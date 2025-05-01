package controller

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/services"
	"authentication/internal/utils"
	"authentication/internal/utils/cdn"
	"authentication/package/response"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetProfile(c *gin.Context)
	UpdateNameUserProfile(c *gin.Context)
	UpdatePhotoUserProfile(c *gin.Context)
	UpdateUserSetting(c *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userController struct {
	UserService services.UserService
	JWTService  utils.JWTService
	ipCdn       string
}

func NewUserController(serviceUser services.UserService, jwtService utils.JWTService, url string) UserController {
	return userController{UserService: serviceUser, JWTService: jwtService, ipCdn: url}
}

func (h userController) GetProfile(c *gin.Context) {
	utils.SendEmail([]string{"morganhero35@gmail.com"}, "Test Email", "This is a test email body")
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

func (h userController) UpdatePhotoUserProfile(context *gin.Context) {
	err := context.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		response.SendResponse(context, 400, "Error parsing form", nil, err.Error())
		return
	}

	file, err := context.FormFile("profile_picture")
	if err != nil {
		response.SendResponse(context, 400, "No image uploaded", nil, "An image is required")
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	var imageMeta out.ImageResponse
	imageMeta, err = cdn.UploadImageToCDN(h.ipCdn, file, token.ClientID, context.GetHeader("Authorization"))
	if err != nil {
		response.SendResponse(context, 500, "Failed to upload image", nil, err.Error())
		return
	}
	if imageMeta.ImageURL == "" {
		response.SendResponse(context, 500, "Failed to upload image", nil, "Image URL is empty")
		return
	}

	user, err := h.UserService.UpdatePhotoUserProfile(imageMeta.ImageURL, token.ClientID)
	if err != nil {
		handleErrorResponse(context, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	handleSuccessResponse(context, http.StatusOK, "Profile updated successfully", user)
}

func (h userController) UpdateUserSetting(c *gin.Context) {
	var req in.UserSettingsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Validation failed",
				"fields": ve.Error(), // or iterate for detailed field errors
			})
		} else {
			handleErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		}
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	errs := h.UserService.UpdateUserSetting(&req, token.ClientID)
	if errs.Message != "" {
		response.SendResponse(c, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	handleSuccessResponse(c, http.StatusOK, "User settings updated successfully", nil)
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
