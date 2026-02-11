package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlan string

const (
	SubscriptionPlanFree   SubscriptionPlan = "free"
	SubscriptionPlanSimple SubscriptionPlan = "simple"
	SubscriptionPlanMedium SubscriptionPlan = "medium"
	SubscriptionPlanUltra  SubscriptionPlan = "ultra"
)

func (p SubscriptionPlan) IsValid() bool {
	switch p {
	case SubscriptionPlanFree, SubscriptionPlanSimple, SubscriptionPlanMedium, SubscriptionPlanUltra:
		return true
	default:
		return false
	}
}

type SubscriptionStatus string

const (
	SubscriptionStatusTrialing             SubscriptionStatus = "trialing"
	SubscriptionStatusActive               SubscriptionStatus = "active"
	SubscriptionStatusPastDue              SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled             SubscriptionStatus = "canceled"
	SubscriptionStatusIncomplete           SubscriptionStatus = "incomplete"
	SubscriptionStatusIncompleteExpired    SubscriptionStatus = "incomplete_expired"
	SubscriptionStatusUnpaid               SubscriptionStatus = "unpaid"
	SubscriptionStatusPaused               SubscriptionStatus = "paused"
	SubscriptionStatusAccessRevokedRefund  SubscriptionStatus = "access_revoked_refund"
	SubscriptionStatusAccessRevokedManual  SubscriptionStatus = "access_revoked_manual"
	SubscriptionStatusAccessRevokedUnknown SubscriptionStatus = "access_revoked_unknown"
)

type Subscription struct {
	gorm.Model
	ID                   uuid.UUID          `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:subscriptions"`
	UserID               uuid.UUID          `json:"user_id" gorm:"type:char(36);not null;uniqueIndex"`
	Plan                 SubscriptionPlan   `json:"plan" gorm:"size:20;not null;default:'free';index"`
	Status               SubscriptionStatus `json:"status" gorm:"size:40;not null;default:'trialing';index"`
	StripeCustomerID     string             `json:"stripe_customer_id" gorm:"size:255;index"`
	StripeSubscriptionID string             `json:"stripe_subscription_id" gorm:"size:255;uniqueIndex"`

	CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty" gorm:"index"`
	CancelAtPeriodEnd  bool       `json:"cancel_at_period_end" gorm:"not null;default:false"`
	CanceledAt         *time.Time `json:"canceled_at,omitempty" gorm:"index"`
	TrialEndsAt        *time.Time `json:"trial_ends_at,omitempty" gorm:"index"`
	AccessRevokedAt    *time.Time `json:"access_revoked_at,omitempty" gorm:"index"`
	AccessRevokeReason *string    `json:"access_revoke_reason,omitempty" gorm:"size:255"`

	User User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
