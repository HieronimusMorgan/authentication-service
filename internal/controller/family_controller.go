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
	UpdateFamily(c *gin.Context)
	AddMemberFamily(c *gin.Context)
	RemoveMemberFamily(c *gin.Context)
	AddFamilyMemberPermission(c *gin.Context)
	RemoveFamilyMemberPermission(c *gin.Context)
	GetListFamilyMemberPermissions(c *gin.Context)
	GetFamilyMembers(c *gin.Context)
}

type familyController struct {
	FamilyService                 services.FamilyService
	FamilyMemberService           services.FamilyMemberService
	FamilyMemberPermissionService services.FamilyMemberPermissionService
	JWTService                    utils.JWTService
}

func NewFamilyController(serviceFamily services.FamilyService, serviceFamilyMember services.FamilyMemberService, serviceFamilyMemberPermission services.FamilyMemberPermissionService, jwtService utils.JWTService) FamilyController {
	return familyController{FamilyService: serviceFamily, FamilyMemberService: serviceFamilyMember, FamilyMemberPermissionService: serviceFamilyMemberPermission, JWTService: jwtService}
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

func (f familyController) UpdateFamily(c *gin.Context) {
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

	err := f.FamilyService.UpdateFamily(&req, token.ClientID)
	if err.Message != "" {
		response.SendResponse(c, err.Code, err.Error, nil, err.Message)
		return
	}

	response.SendResponse(c, http.StatusOK, "Success", nil, "Family updated successfully")
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

	err := f.FamilyMemberService.AddFamilyMember(&req, token.ClientID)
	if err.Message != "" {
		response.SendResponse(c, err.Code, err.Error, nil, err.Message)
		return
	}

	response.SendResponse(c, http.StatusOK, "Success", nil, "Member added successfully")
}

func (f familyController) RemoveMemberFamily(c *gin.Context) {
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

	err := f.FamilyMemberService.RemoveFamilyMember(&req, token.ClientID)
	if err.Message != "" {
		response.SendResponse(c, err.Code, err.Error, nil, err.Message)
		return
	}

	response.SendResponse(c, http.StatusOK, "Success", nil, "Member removed successfully")
}

func (f familyController) AddFamilyMemberPermission(c *gin.Context) {
	var req in.UpdateFamilyMemberPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := f.FamilyMemberPermissionService.AddFamilyMemberPermissions(&req, token.ClientID)
	if err.Message != "" {
		response.SendResponse(c, err.Code, err.Error, nil, err.Message)
		return
	}

	response.SendResponse(c, http.StatusOK, "Success", nil, "Permission changed successfully")
}

func (f familyController) RemoveFamilyMemberPermission(c *gin.Context) {
	var req in.UpdateFamilyMemberPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := f.FamilyMemberPermissionService.RemoveFamilyMemberPermissions(&req, token.ClientID)
	if err.Message != "" {
		response.SendResponse(c, err.Code, err.Error, nil, err.Message)
		return
	}

	response.SendResponse(c, http.StatusOK, "Success", nil, "Permission changed successfully")
}

func (f familyController) GetListFamilyMemberPermissions(c *gin.Context) {
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

	permissions, errs := f.FamilyMemberPermissionService.GetListFamilyMemberPermissions(req, token.ClientID)
	if errs.Message != "" {
		response.SendResponse(c, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	response.SendResponse(c, http.StatusOK, "Success", permissions, "")
}

func (f familyController) GetFamilyMembers(c *gin.Context) {
	familyID, err := utils.ConvertToUint(c.Param("id"))
	if err != nil {
		response.SendResponse(c, 400, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	members, errs := f.FamilyMemberService.GetFamilyMemberByFamilyID(familyID, token.ClientID)
	if errs.Message != "" {
		response.SendResponse(c, errs.Code, errs.Error, nil, errs.Message)
		return
	}

	response.SendResponse(c, http.StatusOK, "Success", members, "")
}
