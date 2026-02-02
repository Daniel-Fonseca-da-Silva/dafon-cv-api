package middleware

import (
	"errors"
	"strings"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StaticTokenMiddleware validates static token from environment variable
func StaticTokenMiddleware(staticToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if static token is configured
		if staticToken == "" {
			transporthttp.HandleValidationError(c, errors.New("static token not configured"))
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			transporthttp.HandleValidationError(c, errors.New("authorization header required"))
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			transporthttp.HandleValidationError(c, errors.New("invalid authorization header format"))
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token against the static token
		if tokenString != staticToken {
			transporthttp.HandleValidationError(c, errors.New("invalid token"))
			c.Abort()
			return
		}

		// Token is valid, proceed to the next handler
		c.Next()
	}
}

// AdminMiddleware validates X-User-ID header and ensures the user exists and has Admin=true
func AdminMiddleware(userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "X-User-ID header required"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid user ID format"})
			return
		}

		user, err := userRepo.GetByID(c.Request.Context(), userID)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(404, gin.H{"error": "user not found"})
			return
		}

		if !user.Admin {
			c.AbortWithStatusJSON(403, gin.H{"error": "admin access required"})
			return
		}

		c.Next()
	}
}
