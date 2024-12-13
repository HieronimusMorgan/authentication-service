package utils

import (
	"Authentication/internal/models"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var internalSecretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(user models.User) (models.TokenDetails, error) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.New().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.New().String()

	claims := jwt.MapClaims{
		"authorized":  true,
		"access_uuid": td.AccessUUID,
		"user_id":     user.UserID,
		"client_id":   user.ClientID,
		"role_id":     user.RoleID,
		"exp":         td.AtExpires,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	td.AccessToken, err = at.SignedString(jwtSecret)
	if err != nil {
		return models.TokenDetails{}, err
	}

	rtClaims := jwt.MapClaims{
		"refresh_uuid": td.RefreshUUID,
		"user_id":      user.UserID,
		"exp":          td.RtExpires,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(jwtSecret)
	if err != nil {
		return models.TokenDetails{}, err
	}

	return *td, nil
}

// ValidateToken validates the JWT and extracts claims
func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Validate token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			if time.Now().After(expTime) {
				return nil, errors.New("token is expired")
			}
		} else {
			return nil, errors.New("missing exp claim")
		}
		return token, nil
	}

	return nil, errors.New("invalid token claims")
}

// TokenClaims represents the extracted claims from the JWT
type TokenClaims struct {
	Authorized bool   `json:"authorized"`
	AccessUUID string `json:"access_uuid"`
	UserID     uint   `json:"user_id"`
	ClientID   string `json:"client_id"`
	Role       string `json:"role"`
}

// ExtractClaims validates a JWT and extracts claims
func ExtractClaims(tokenString string) (*TokenClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Populate TokenClaims struct
	tc := &TokenClaims{}

	// Extract `authorized` claim
	if authorized, ok := claims["authorized"].(bool); ok {
		tc.Authorized = authorized
	}

	// Extract `access_uuid` claim
	if accessUUID, ok := claims["access_uuid"].(string); ok {
		tc.AccessUUID = accessUUID
	}

	// Extract `user_id` claim
	if userID, ok := claims["user_id"].(float64); ok { // JWT numbers are float64
		tc.UserID = uint(userID)
	}

	// Extract `client_id` claim
	if clientID, ok := claims["client_id"].(string); ok { // JWT numbers are float64
		tc.ClientID = clientID
	}

	// Extract `role` claim
	if role, ok := claims["role"].(string); ok {
		tc.Role = role
	}

	return tc, nil
}
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}

func GetClientIDFromToken(tokenString string) (interface{}, error) {
	token, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	clientID, exists := claims["client_id"]
	if !exists {
		return nil, errors.New("client_id not found in token claims")
	}

	return clientID, nil
}

func GetRoleIDFromToken(tokenString string) (interface{}, error) {
	token, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	clientID, exists := claims["role_id"]
	if !exists {
		return nil, errors.New("client_id not found in token claims")
	}

	return clientID, nil
}

func GetUserIDFromToken(tokenString string) (uint, error) {
	token, err := ValidateJWT(tokenString)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	// Extract user_id from claims
	userIDFloat, exists := claims["user_id"]
	if !exists {
		return 0, errors.New("user_id not found in token claims")
	}

	// Convert user_id to uint
	userIDFloat64, ok := userIDFloat.(float64)
	if !ok {
		return 0, errors.New("user_id is not a valid number")
	}

	userID := uint(userIDFloat64)
	return userID, nil
}

func GetUUIDFromToken(tokenString string) (interface{}, error) {
	token, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	clientID, exists := claims["uuid_key"]
	if !exists {
		return nil, errors.New("user_id not found in token claims")
	}

	return clientID, nil
}

func GetExpFromToken(tokenString string) (interface{}, error) {
	token, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	clientID, exists := claims["exp"]
	if !exists {
		return nil, errors.New("user_id not found in token claims")
	}

	return clientID, nil
}

func Float64ToUint(value interface{}) (uint, error) {
	floatVal, ok := value.(float64)
	if !ok {
		return 0, errors.New("value is not a float64")
	}

	if floatVal < 0 {
		return 0, errors.New("cannot convert negative float64 to uint")
	}

	return uint(floatVal), nil
}

// Secret key for signing the tokens

type InternalClaims struct {
	Service string `json:"service"` // Service name
	jwt.RegisteredClaims
}

func GenerateInternalToken(serviceName string) (string, error) {
	// Define claims
	claims := InternalClaims{
		Service: serviceName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth-service",                                      // Issuing service
			Subject:   "internal-communication",                            // Purpose
			Audience:  []string{strings.ToLower(serviceName) + "-service"}, // Target service(s)
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),       // Token expiration (1 hour)
			IssuedAt:  jwt.NewNumericDate(time.Now()),                      // Issued time
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	return token.SignedString(internalSecretKey)
}

func ValidateInternalToken(tokenString string) (*InternalClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &InternalClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return internalSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token claims
	if claims, ok := token.Claims.(*InternalClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
