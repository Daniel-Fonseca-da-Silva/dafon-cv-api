package middleware

import (
	"errors"
	"strings"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware(jwtConfig *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.HandleValidationError(c, errors.New("Authorization header required"))
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.HandleValidationError(c, errors.New("invalid authorization header format"))
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token with claims
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtConfig.SecretKey), nil
		})

		if err != nil {
			// Check if the error is specifically about expiration
			if errors.Is(err, jwt.ErrTokenExpired) {
				utils.HandleValidationError(c, errors.New("token expired"))
				c.Abort()
				return
			}
			utils.HandleValidationError(c, errors.New("invalid token"))
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			utils.HandleValidationError(c, errors.New("invalid token"))
			c.Abort()
			return
		}

		// Extract user ID from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			utils.HandleValidationError(c, errors.New("invalid user ID in token"))
			c.Abort()
			return
		}

		// Set user ID in context for later use
		c.Set("user_id", userID)
		c.Next()
	}
}
