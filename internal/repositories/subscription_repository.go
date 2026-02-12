package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SubscriptionRepository defines the interface for subscription data operations.
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error)
	GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*models.Subscription, error)
	GetByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*models.Subscription, error)
	Save(ctx context.Context, subscription *models.Subscription) error
}

type subscriptionRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewSubscriptionRepository(db *gorm.DB, logger *zap.Logger) SubscriptionRepository {
	return &subscriptionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	if err := r.db.WithContext(ctx).Create(subscription).Error; err != nil {
		r.logger.Error(
			"Failed to create subscription",
			zap.Error(err),
			zap.String("user_id", subscription.UserID.String()),
		)
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error(
			"Failed to get subscription by user ID",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return nil, fmt.Errorf("failed to get subscription by user ID %s: %w", userID.String(), err)
	}
	return &subscription, nil
}

func (r *subscriptionRepository) GetByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.WithContext(ctx).Where("stripe_subscription_id = ?", stripeSubscriptionID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error(
			"Failed to get subscription by Stripe subscription ID",
			zap.Error(err),
			zap.String("stripe_subscription_id", stripeSubscriptionID),
		)
		return nil, fmt.Errorf("failed to get subscription by Stripe subscription ID %s: %w", stripeSubscriptionID, err)
	}
	return &subscription, nil
}

func (r *subscriptionRepository) GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.WithContext(ctx).Where("stripe_customer_id = ?", stripeCustomerID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error(
			"Failed to get subscription by Stripe customer ID",
			zap.Error(err),
			zap.String("stripe_customer_id", stripeCustomerID),
		)
		return nil, fmt.Errorf("failed to get subscription by Stripe customer ID %s: %w", stripeCustomerID, err)
	}
	return &subscription, nil
}

func (r *subscriptionRepository) Save(ctx context.Context, subscription *models.Subscription) error {
	if err := r.db.WithContext(ctx).Save(subscription).Error; err != nil {
		r.logger.Error(
			"Failed to save subscription",
			zap.Error(err),
			zap.String("user_id", subscription.UserID.String()),
		)
		return fmt.Errorf("failed to save subscription: %w", err)
	}
	return nil
}
