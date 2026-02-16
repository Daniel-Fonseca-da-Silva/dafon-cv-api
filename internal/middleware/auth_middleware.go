package middleware

import (
	"net/http"
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
			transporthttp.HandleError(c, http.StatusInternalServerError, "static token not configured")
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			transporthttp.HandleError(c, http.StatusUnauthorized, "authorization header required")
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			transporthttp.HandleError(c, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			transporthttp.HandleError(c, http.StatusUnauthorized, "authorization token required")
			return
		}

		// Validate the token against the static token
		if tokenString != staticToken {
			transporthttp.HandleError(c, http.StatusUnauthorized, "invalid token")
			return
		}

		// Token is valid, proceed to the next handler
		c.Next()
	}
}

// StaticTokenHeaderName is the header used for static token in admin routes (so Authorization can carry session token).
const StaticTokenHeaderName = "X-Static-Token"

// StaticTokenHeaderMiddleware validates static token from X-Static-Token header.
// Used for admin routes so that Authorization can carry the session token.
// Ensures the request comes from a trusted client (e.g. Next.js) that holds the static token.
func StaticTokenHeaderMiddleware(staticToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if staticToken == "" {
			transporthttp.HandleError(c, http.StatusInternalServerError, "static token not configured")
			return
		}

		headerToken := c.GetHeader(StaticTokenHeaderName)
		if headerToken == "" {
			transporthttp.HandleError(c, http.StatusUnauthorized, "X-Static-Token header required")
			return
		}

		if headerToken != staticToken {
			transporthttp.HandleError(c, http.StatusUnauthorized, "invalid static token")
			return
		}

		c.Next()
	}
}

// SessionMiddleware validates a per-user session token (magic link login) and
// sets the authenticated user id in Gin context under key "user_id".
func SessionMiddleware(sessionRepo repositories.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			transporthttp.HandleError(c, http.StatusUnauthorized, "authorization header required")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			transporthttp.HandleError(c, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			transporthttp.HandleError(c, http.StatusUnauthorized, "authorization token required")
			return
		}

		if sessionRepo == nil {
			transporthttp.HandleError(c, http.StatusInternalServerError, "session repository not configured")
			return
		}

		session, err := sessionRepo.GetByToken(tokenString)
		if err != nil {
			transporthttp.HandleError(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if session == nil {
			transporthttp.HandleError(c, http.StatusUnauthorized, "invalid or expired session")
			return
		}

		c.Set("user_id", session.UserID)
		c.Next()
	}
}

// AdminMiddleware validates X-User-ID header and ensures the user exists and has Admin=true
func AdminMiddleware(userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if userRepo == nil {
			transporthttp.HandleError(c, http.StatusInternalServerError, "user repository not configured")
			return
		}

		if ctxUserID, ok := c.Get("user_id"); ok {
			userID, ok := ctxUserID.(uuid.UUID)
			if !ok {
				transporthttp.HandleError(c, http.StatusInternalServerError, "invalid user id in request context")
				return
			}

			ensureAdminOrAbort(c, userRepo, userID)
			return
		}

		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			transporthttp.HandleError(c, http.StatusUnauthorized, "X-User-ID header required")
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			transporthttp.HandleError(c, http.StatusBadRequest, "invalid user ID format")
			return
		}

		ensureAdminOrAbort(c, userRepo, userID)
	}
}

func ensureAdminOrAbort(c *gin.Context, userRepo repositories.UserRepository, userID uuid.UUID) {
	user, err := userRepo.GetByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		transporthttp.HandleError(c, http.StatusNotFound, "user not found")
		return
	}

	if !user.Admin {
		transporthttp.HandleError(c, http.StatusForbidden, "admin access required")
		return
	}

	c.Next()
}
