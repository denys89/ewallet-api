package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/denys89/ewallet-api/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const UserIDKey = "user_id"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondWithError(c, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Extract the Bearer token
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			respondWithError(c, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// Return the secret key for validation
			return []byte(config.Get().JWTSecret), nil
		})

		if err != nil || !token.Valid {
			respondWithError(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Extract and validate claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondWithError(c, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Extract user_id from claims
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			respondWithError(c, http.StatusUnauthorized, "Invalid user ID in token")
			return
		}
		log.Printf("userIDStr: %+v\n", userIDStr)

		// Validate user_id format (UUID)
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			respondWithError(c, http.StatusUnauthorized, "Invalid user ID in token")
			return
		}

		// Store the user ID in the context for later use

		c.Set(UserIDKey, userID)

		// Proceed to the next middleware or handler
		c.Next()
	}
}

// Helper function for responding with an error
func respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
