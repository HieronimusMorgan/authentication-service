package controller

import (
	"authentication/internal/dto/in"
	"authentication/internal/services"
	"authentication/internal/utils"
	"authentication/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FamilyController interface {
	CreateFamily(c *gin.Context)
	AddMemberFamily(c *gin.Context)
}

type familyController struct {
	FamilyService services.FamilyService
	JWTService    utils.JWTService
}

func NewFamilyController(serviceFamily services.FamilyService, jwtService utils.JWTService) FamilyController {
	return familyController{FamilyService: serviceFamily, JWTService: jwtService}
}

func (f familyController) CreateFamily(c *gin.Context) {
	var req in.FamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := f.FamilyService.CreateFamily(&req, token.ClientID)
	if err.Message != "" {
		response.SendResponse(c, err.Code, err.Error, nil, err.Message)
		return
	}

	response.SendResponse(c, http.StatusOK, "Success", nil, "Family created successfully")
}

func (f familyController) AddMemberFamily(c *gin.Context) {
	var req in.FamilyMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := f.FamilyService.CreateFamilyMember(&req, token.ClientID)
	if err.Message != "" {
		response.SendResponse(c, err.Code, err.Error, nil, err.Message)
		return
	}

	response.SendResponse(c, http.StatusOK, "Success", nil, "Member added successfully")
}
