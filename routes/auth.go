package routes

import (
	"net/http"
	"time"

	"github.com/denys89/ewallet-api/config"
	"github.com/denys89/ewallet-api/models"
	"github.com/denys89/ewallet-api/repositories"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Pin         string `json:"pin" binding:"required,len=6"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.Pin), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash PIN"})
		return
	}

	user := models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Pin:         string(hashedPin),
		Balance:     0,
	}

	userRepo := repositories.NewUserRepository(config.DB)
	if err := userRepo.Create(&user); err != nil {
		if err == repositories.ErrPhoneNumberExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone Number already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"user_id":      user.ID,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"phone_number": user.PhoneNumber,
			"address":      user.Address,
			"created_date": user.CreatedAt.Format("2006-1-2 15:04:05"),
		},
	})
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Pin         string `json:"pin" binding:"required,len=6"`
}

// generateTokens creates both access and refresh tokens
func generateTokens(userID string, phoneNumber string) (string, string, error) {
	cfg := config.Get()

	// Generate access token
	accessTokenExp := time.Now().Add(cfg.JWTExpirationHours * time.Hour)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"phone":   phoneNumber,
		"exp":     accessTokenExp.Unix(),
		"iat":     time.Now().Unix(),
		"type":    "access",
	})

	accessTokenString, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshTokenExp := time.Now().Add(cfg.RefreshTokenExpirationDays * 24 * time.Hour)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"phone":   phoneNumber,
		"exp":     refreshTokenExp.Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.RefreshTokenSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repositories.NewUserRepository(config.DB)
	user, err := userRepo.FindByPhoneNumber(req.PhoneNumber)
	if err != nil {
		if err == repositories.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Phone Number and PIN doesn't match"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Phone Number and PIN doesn't match"})
		return
	}

	accessToken, refreshToken, err := generateTokens(user.ID.String(), user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    int(config.Get().JWTExpirationHours * 3600),
		},
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken handles token refresh requests
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse and validate the refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Get().RefreshTokenSecret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Verify token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
		return
	}

	// Get user information from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	phoneNumber, ok := claims["phone"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number in token"})
		return
	}

	// Generate new tokens
	accessToken, refreshToken, err := generateTokens(userID, phoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    int(config.Get().JWTExpirationHours * 3600),
		},
	})
}
