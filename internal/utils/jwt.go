package utils

import (
	"authentication/internal/models"
	"authentication/package/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type JWTService interface {
	GenerateToken(user models.Users) (models.TokenDetails, error)
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
	ValidateTokenAdmin(tokenString string) (*jwt.MapClaims, error)
	ExtractClaims(tokenString string) (*TokenClaims, error)
	GenerateInternalToken(serviceName string) (string, error)
	ValidateInternalToken(tokenString string) (*InternalClaims, error)
}

type jwtService struct {
	SecretKey         []byte
	InternalSecretKey []byte
}

// NewJWTService initializes the JWT service
func NewJWTService(jwtSecret string) JWTService {
	return jwtService{
		SecretKey:         []byte(jwtSecret),
		InternalSecretKey: []byte(jwtSecret),
	}
}

// GenerateToken generates a new JWT token
func (j jwtService) GenerateToken(user models.Users) (models.TokenDetails, error) {
	td := &models.TokenDetails{
		AtExpires:   time.Now().Add(time.Hour * 24).Unix(),
		AccessUUID:  uuid.New().String(),
		RtExpires:   time.Now().Add(time.Hour * 24 * 7).Unix(),
		RefreshUUID: GenerateClientID(),
	}

	claims := jwt.MapClaims{
		"authorized":  true,
		"access_uuid": td.AccessUUID,
		"user_id":     user.UserID,
		"client_id":   user.ClientID,
		"role_id":     user.RoleID,
		"role":        user.Role.Name,
		"exp":         td.AtExpires,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	td.AccessToken, err = at.SignedString(j.SecretKey)
	if err != nil {
		return models.TokenDetails{}, err
	}

	td.RefreshToken = GenerateClientID()

	return *td, nil
}

// ValidateToken validates a JWT token and extracts claims
func (j jwtService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
	}

	return &claims, nil
}

// ValidateTokenAdmin checks if a user is an admin based on JWT claims
func (j jwtService) ValidateTokenAdmin(tokenString string) (*jwt.MapClaims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if role, ok := (*claims)["role"].(string); ok {
		if strings.EqualFold(role, "Admin") || strings.EqualFold(role, "Super Admin") {
			return claims, nil
		}
		return nil, errors.New("user is not an Admin")
	}

	return nil, errors.New("role not found in token claims")
}

// ExtractClaims extracts claims from a JWT token
func (j jwtService) ExtractClaims(tokenString string) (*TokenClaims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	tc := &TokenClaims{}

	if authorized, ok := (*claims)["authorized"].(bool); ok {
		tc.Authorized = authorized
	}

	if accessUUID, ok := (*claims)["access_uuid"].(string); ok {
		tc.AccessUUID = accessUUID
	}

	if exp, ok := (*claims)["exp"].(float64); ok {
		tc.Exp = int64(exp)
	}

	if userID, ok := (*claims)["user_id"].(float64); ok {
		tc.UserID = uint(userID)
	}

	if clientID, ok := (*claims)["client_id"].(string); ok {
		tc.ClientID = clientID
	}

	if role, ok := (*claims)["role_id"].(float64); ok {
		tc.RoleID = uint(role)
	}

	return tc, nil
}

// GenerateInternalToken creates an internal JWT for service-to-service communication
func (j jwtService) GenerateInternalToken(serviceName string) (string, error) {
	claims := InternalClaims{
		Service: serviceName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth-service",
			Subject:   "internal-communication",
			Audience:  []string{strings.ToLower(serviceName) + "-service"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.InternalSecretKey)
}

// ValidateInternalToken verifies an internal JWT token
func (j jwtService) ValidateInternalToken(tokenString string) (*InternalClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &InternalClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.InternalSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*InternalClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// TokenClaims represents the claims extracted from a JWT token
type TokenClaims struct {
	Authorized bool   `json:"authorized"`
	AccessUUID string `json:"access_uuid"`
	UserID     uint   `json:"user_id"`
	ClientID   string `json:"client_id"`
	RoleID     uint   `json:"role_id"`
	Exp        int64  `json:"exp"`
}

// InternalClaims represents the claims used for service-to-service authentication
type InternalClaims struct {
	Service string `json:"service"`
	jwt.RegisteredClaims
}

// ExtractTokenClaims extracts token claims from the Gin context
func ExtractTokenClaims(c *gin.Context) (*TokenClaims, bool) {
	tokenData, exists := c.Get("token")
	if !exists {
		response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "Token not found")
		return nil, false
	}

	tokenClaims, ok := tokenData.(*TokenClaims)
	if !ok || tokenClaims == nil {
		response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "Invalid token format")
		return nil, false
	}

	return tokenClaims, true
}
