package handlers

import (
	"errors"
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	subscriptionUseCase usecases.SubscriptionUseCase
	logger              *zap.Logger
}

func NewSubscriptionHandler(subscriptionUseCase usecases.SubscriptionUseCase, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionUseCase: subscriptionUseCase,
		logger:              logger,
	}
}

func (h *SubscriptionHandler) GetMySubscription(c *gin.Context) {
	userID, ok := h.getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.subscriptionUseCase.GetMySubscription(c.Request.Context(), userID)
	if err != nil {
		h.abortWithInternalServerError(c, "get my subscription", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) CreateCheckoutSession(c *gin.Context) {
	userID, ok := h.getUserIDFromContext(c)
	if !ok {
		return
	}

	var req dto.CreateCheckoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	resp, err := h.subscriptionUseCase.CreateCheckoutSession(c.Request.Context(), userID, &req)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) CreatePortalSession(c *gin.Context) {
	userID, ok := h.getUserIDFromContext(c)
	if !ok {
		return
	}

	var req dto.CreatePortalSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	resp, err := h.subscriptionUseCase.CreatePortalSession(c.Request.Context(), userID, &req)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) CancelMySubscription(c *gin.Context) {
	userID, ok := h.getUserIDFromContext(c)
	if !ok {
		return
	}

	var req dto.CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	if err := h.subscriptionUseCase.CancelMySubscription(c.Request.Context(), userID, req.AtPeriodEnd); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancellation requested"})
}

func (h *SubscriptionHandler) StripeWebhook(c *gin.Context) {
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		if h.logger != nil {
			h.logger.Warn("Stripe webhook rejected: missing Stripe-Signature header",
				zap.String("path", c.FullPath()),
			)
		}
		transporthttp.HandleValidationError(c, errors.New("stripe signature header required"))
		return
	}

	payload, err := c.GetRawData()
	if err != nil {
		h.abortWithInternalServerError(c, "read stripe webhook body", err)
		return
	}

	if err := h.subscriptionUseCase.HandleStripeWebhook(c.Request.Context(), payload, signature); err != nil {
		if h.logger != nil {
			h.logger.Error("Stripe webhook handling failed",
				zap.String("path", c.FullPath()),
				zap.Int("payload_size", len(payload)),
				zap.Int("signature_size", len(signature)),
				zap.Error(err),
			)
		}
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *SubscriptionHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	ctxUserID, ok := c.Get("user_id")
	if !ok {
		transporthttp.HandleValidationError(c, errors.New("user not authenticated"))
		return uuid.Nil, false
	}

	userID, ok := ctxUserID.(uuid.UUID)
	if !ok {
		transporthttp.HandleValidationError(c, errors.New("invalid user id in request context"))
		return uuid.Nil, false
	}

	return userID, true
}

func (h *SubscriptionHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Subscription handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
