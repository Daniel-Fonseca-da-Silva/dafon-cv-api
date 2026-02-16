package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCheckoutSessionRequest struct {
	Plan       string `json:"plan" validate:"required,oneof=simple medium ultra"`
	SuccessURL string `json:"success_url" validate:"required,url"`
	CancelURL  string `json:"cancel_url" validate:"required,url"`
}

type CreateCheckoutSessionResponse struct {
	SessionID   string `json:"session_id"`
	CheckoutURL string `json:"checkout_url"`
}

type CreatePortalSessionRequest struct {
	ReturnURL string `json:"return_url" validate:"required,url"`
}

type CreatePortalSessionResponse struct {
	PortalURL string `json:"portal_url"`
}

type CancelSubscriptionRequest struct {
	AtPeriodEnd bool `json:"at_period_end"`
}

type SubscriptionResponse struct {
	UserID            uuid.UUID  `json:"user_id"`
	Plan              string     `json:"plan"`
	Status            string     `json:"status"`
	StripeCustomerID  *string    `json:"stripe_customer_id,omitempty"`
	CurrentPeriodEnd  *time.Time `json:"current_period_end,omitempty"`
	CancelAtPeriodEnd bool       `json:"cancel_at_period_end"`
	CanceledAt        *time.Time `json:"canceled_at,omitempty"`
	TrialEndsAt       *time.Time `json:"trial_ends_at,omitempty"`
	AccessRevokedAt   *time.Time `json:"access_revoked_at,omitempty"`
}
