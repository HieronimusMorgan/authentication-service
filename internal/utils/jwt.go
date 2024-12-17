package utils

import (
	"authentication/internal/models"
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

func GenerateToken(user models.Users) (models.TokenDetails, error) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 24).Unix()
	td.AccessUUID = uuid.New().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = GenerateClientID()

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
	td.AccessToken, err = at.SignedString(jwtSecret)
	if err != nil {
		return models.TokenDetails{}, err
	}

	td.RefreshToken = GenerateClientID()

	return *td, nil
}

func ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims and validate
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
	}

	return &claims, nil
}

func ValidateTokenAdmin(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
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

	if role, ok := claims["role"].(string); ok {
		if strings.EqualFold(role, "Admin") || strings.EqualFold(role, "Super Admin") {
			return &claims, nil
		}
		return nil, errors.New("user is not an Admin")
	}

	return nil, errors.New("role not found in token claims")
}

func ExtractClaims(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	tc := &TokenClaims{}

	if authorized, ok := claims["authorized"].(bool); ok {
		tc.Authorized = authorized
	}

	if accessUUID, ok := claims["access_uuid"].(string); ok {
		tc.AccessUUID = accessUUID
	}

	if exp, ok := claims["exp"].(float64); ok {
		tc.Exp = int64(exp)
	}

	if userID, ok := claims["user_id"].(float64); ok { // JWT numbers are float64
		tc.UserID = uint(userID)
	}

	if clientID, ok := claims["client_id"].(string); ok { // JWT numbers are float64
		tc.ClientID = clientID
	}

	if role, ok := claims["role_id"].(float64); ok {
		tc.RoleID = uint(role)
	}

	return tc, nil
}

func GenerateInternalToken(serviceName string) (string, error) {
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
	token, err := jwt.ParseWithClaims(tokenString, &InternalClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return internalSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*InternalClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

type TokenClaims struct {
	Authorized bool   `json:"authorized"`
	AccessUUID string `json:"access_uuid"`
	UserID     uint   `json:"user_id"`
	ClientID   string `json:"client_id"`
	RoleID     uint   `json:"role_id"`
	Exp        int64  `json:"exp"`
}

type InternalClaims struct {
	Service string `json:"service"` // Service name
	jwt.RegisteredClaims
}
