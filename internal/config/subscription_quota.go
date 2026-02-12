package config

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
)

// Env keys for subscription quota (monthly AI requests per plan).
const (
	envQuotaFreeMonthly   = "SUBSCRIPTION_QUOTA_FREE_MONTHLY"
	envQuotaSimpleMonthly = "SUBSCRIPTION_QUOTA_SIMPLE_MONTHLY"
	envQuotaMediumMonthly = "SUBSCRIPTION_QUOTA_MEDIUM_MONTHLY"
	envQuotaUltraMonthly  = "SUBSCRIPTION_QUOTA_ULTRA_MONTHLY"
)

// PlanQuota holds the monthly request limit for a subscription plan.
type PlanQuota struct {
	MonthlyRequests int64
}

// DefaultAIQuotaByPlan returns the monthly request quota per subscription plan,
// read from env (SUBSCRIPTION_QUOTA_*_MONTHLY) with fallback to defaults.
func DefaultAIQuotaByPlan() map[models.SubscriptionPlan]PlanQuota {
	return map[models.SubscriptionPlan]PlanQuota{
		models.SubscriptionPlanFree:   {MonthlyRequests: int64(ParseIntEnv(envQuotaFreeMonthly, 10))},
		models.SubscriptionPlanSimple: {MonthlyRequests: int64(ParseIntEnv(envQuotaSimpleMonthly, 30))},
		models.SubscriptionPlanMedium: {MonthlyRequests: int64(ParseIntEnv(envQuotaMediumMonthly, 100))},
		models.SubscriptionPlanUltra:  {MonthlyRequests: int64(ParseIntEnv(envQuotaUltraMonthly, -1))},
	}
}
