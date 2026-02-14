package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupSubscriptionRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config, authMiddleware gin.HandlerFunc) {
	userRepo := repositories.NewUserRepository(db, logger)
	subscriptionRepo := repositories.NewSubscriptionRepository(db, logger)
	subscriptionUseCase := usecases.NewSubscriptionUseCase(subscriptionRepo, userRepo, cfg.Stripe, logger)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionUseCase, logger)

	// Stripe webhook should not be protected by static token.
	router.POST("/api/v1/subscriptions/webhook", subscriptionHandler.StripeWebhook)

	protected := router.Group("/api/v1/subscriptions", authMiddleware)
	{
		protected.GET("/me", subscriptionHandler.GetMySubscription)
		protected.POST("/checkout", subscriptionHandler.CreateCheckoutSession)
		protected.POST("/portal", subscriptionHandler.CreatePortalSession)
	}
}
