package middleware

import (
	"errors"
	"strings"

	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/gin-gonic/gin"
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
