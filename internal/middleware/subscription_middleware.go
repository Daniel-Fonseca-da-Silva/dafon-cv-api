package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func RequireSubscriptionPlan(
	subscriptionUseCase usecases.SubscriptionUseCase,
	redisClient *redis.Client,
	quotaByPlan map[models.SubscriptionPlan]config.PlanQuota,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := parseUserIDHeader(c)
		if !ok {
			return
		}

		plan, err := subscriptionUseCase.GetEntitlement(c.Request.Context(), userID)
		if err != nil {
			transporthttp.HandleError(c, http.StatusInternalServerError, err.Error())
			return
		}

		quota, hasQuota := quotaByPlan[plan]
		if !hasQuota {
			transporthttp.HandleError(c, http.StatusInternalServerError, fmt.Sprintf("no quota configured for plan %q", plan))
			return
		}

		if quota.MonthlyRequests < 0 {
			c.Next()
			return
		}

		if redisClient == nil {
			transporthttp.HandleError(c, http.StatusInternalServerError, "redis client not configured")
			return
		}

		key := buildMonthlyUsageKey(userID, "ai_requests")
		count, err := redisClient.Incr(c.Request.Context(), key).Result()
		if err != nil {
			transporthttp.HandleError(c, http.StatusInternalServerError, fmt.Sprintf("increment usage counter: %v", err))
			return
		}

		expireAt := firstMomentOfNextMonth(time.Now().UTC())
		if err := redisClient.ExpireAt(c.Request.Context(), key, expireAt).Err(); err != nil {
			transporthttp.HandleError(c, http.StatusInternalServerError, fmt.Sprintf("set usage counter expiry: %v", err))
			return
		}

		if quota.MonthlyRequests == 0 || count > quota.MonthlyRequests {
			transporthttp.HandleError(c, http.StatusPaymentRequired, "plan limit exceeded")
			return
		}

		c.Next()
	}
}

func parseUserIDHeader(c *gin.Context) (uuid.UUID, bool) {
	if ctxUserID, ok := c.Get("user_id"); ok {
		if userID, ok := ctxUserID.(uuid.UUID); ok {
			return userID, true
		}
		transporthttp.HandleError(c, http.StatusInternalServerError, "invalid user id in request context")
		return uuid.Nil, false
	}

	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		transporthttp.HandleError(c, http.StatusUnauthorized, "user not authenticated")
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		transporthttp.HandleError(c, http.StatusBadRequest, "invalid user ID format")
		return uuid.Nil, false
	}

	return userID, true
}

func buildMonthlyUsageKey(userID uuid.UUID, feature string) string {
	now := time.Now().UTC()
	return fmt.Sprintf("usage:%s:%04d-%02d:%s", userID.String(), now.Year(), int(now.Month()), feature)
}

func firstMomentOfNextMonth(t time.Time) time.Time {
	firstOfThisMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	return firstOfThisMonth.AddDate(0, 1, 0)
}
