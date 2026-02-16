package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	billingportalsession "github.com/stripe/stripe-go/v81/billingportal/session"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	stripesub "github.com/stripe/stripe-go/v81/subscription"
	"github.com/stripe/stripe-go/v81/webhook"
	"go.uber.org/zap"
)

type SubscriptionUseCase interface {
	GetMySubscription(ctx context.Context, userID uuid.UUID) (*dto.SubscriptionResponse, error)
	GetEntitlement(ctx context.Context, userID uuid.UUID) (models.SubscriptionPlan, error)
	CreateCheckoutSession(ctx context.Context, userID uuid.UUID, req *dto.CreateCheckoutSessionRequest) (*dto.CreateCheckoutSessionResponse, error)
	CreatePortalSession(ctx context.Context, userID uuid.UUID, req *dto.CreatePortalSessionRequest) (*dto.CreatePortalSessionResponse, error)
	HandleStripeWebhook(ctx context.Context, payload []byte, signature string) error
}

type subscriptionUseCase struct {
	subscriptionRepo repositories.SubscriptionRepository
	userRepo         repositories.UserRepository
	logger           *zap.Logger
	stripeCfg        config.StripeConfig
	now              func() time.Time
}

func NewSubscriptionUseCase(
	subscriptionRepo repositories.SubscriptionRepository,
	userRepo repositories.UserRepository,
	stripeCfg config.StripeConfig,
	logger *zap.Logger,
) SubscriptionUseCase {
	return &subscriptionUseCase{
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
		logger:           logger,
		stripeCfg:        stripeCfg,
		now:              time.Now,
	}
}

func (uc *subscriptionUseCase) GetMySubscription(ctx context.Context, userID uuid.UUID) (*dto.SubscriptionResponse, error) {
	subscription, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		return &dto.SubscriptionResponse{
			UserID: userID,
			Plan:   string(models.SubscriptionPlanFree),
			Status: "none",
		}, nil
	}

	planToReturn := subscription.Plan
	if active, _ := isSubscriptionActive(uc.now(), subscription); !active {
		planToReturn = models.SubscriptionPlanFree
	} else if !planToReturn.IsValid() {
		planToReturn = models.SubscriptionPlanFree
	}

	return &dto.SubscriptionResponse{
		UserID:            subscription.UserID,
		Plan:              string(planToReturn),
		Status:            string(subscription.Status),
		StripeCustomerID:  subscription.StripeCustomerID,
		CurrentPeriodEnd:  subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd: subscription.CancelAtPeriodEnd,
		CanceledAt:        subscription.CanceledAt,
		TrialEndsAt:       subscription.TrialEndsAt,
		AccessRevokedAt:   subscription.AccessRevokedAt,
	}, nil
}

func (uc *subscriptionUseCase) GetEntitlement(ctx context.Context, userID uuid.UUID) (models.SubscriptionPlan, error) {
	subscription, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return models.SubscriptionPlanFree, err
	}
	if subscription == nil {
		return models.SubscriptionPlanFree, nil
	}

	active, reason := isSubscriptionActive(uc.now(), subscription)
	if !active {
		uc.logger.Info(
			"Subscription not active, falling back to free plan",
			zap.String("user_id", userID.String()),
			zap.String("reason", reason),
			zap.String("status", string(subscription.Status)),
			zap.String("plan", string(subscription.Plan)),
		)
		return models.SubscriptionPlanFree, nil
	}

	if !subscription.Plan.IsValid() {
		return models.SubscriptionPlanFree, fmt.Errorf("invalid subscription plan for user %s: %s", userID.String(), subscription.Plan)
	}

	return subscription.Plan, nil
}

func (uc *subscriptionUseCase) CreateCheckoutSession(ctx context.Context, userID uuid.UUID, req *dto.CreateCheckoutSessionRequest) (*dto.CreateCheckoutSessionResponse, error) {
	if uc.stripeCfg.SecretKey == "" {
		return nil, errors.New("stripe secret key not configured")
	}

	plan := models.SubscriptionPlan(req.Plan)
	if !plan.IsValid() || plan == models.SubscriptionPlanFree {
		return nil, fmt.Errorf("invalid plan: %s", req.Plan)
	}

	priceID, err := uc.priceIDForPlan(plan)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user %s: %w", userID.String(), err)
	}

	stripe.Key = uc.stripeCfg.SecretKey

	subscription, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	stripeCustomerID := ""
	if subscription != nil && subscription.StripeCustomerID != nil && *subscription.StripeCustomerID != "" {
		stripeCustomerID = *subscription.StripeCustomerID
	} else {
		params := &stripe.CustomerParams{
			Email: stripe.String(user.Email),
			Name:  stripe.String(user.Name),
		}
		params.AddMetadata("user_id", userID.String())

		stripeCustomer, err := customer.New(params)
		if err != nil {
			return nil, fmt.Errorf("create stripe customer: %w", err)
		}
		stripeCustomerID = stripeCustomer.ID
	}

	if subscription == nil {
		subscription = &models.Subscription{
			UserID:           userID,
			Plan:             plan,
			Status:           models.SubscriptionStatusIncomplete,
			StripeCustomerID: &stripeCustomerID,
		}
		if err := uc.subscriptionRepo.Create(ctx, subscription); err != nil {
			return nil, err
		}
	} else {
		subscription.Plan = plan
		subscription.Status = models.SubscriptionStatusIncomplete
		subscription.StripeCustomerID = &stripeCustomerID
		if err := uc.subscriptionRepo.Save(ctx, subscription); err != nil {
			return nil, err
		}
	}

	checkoutParams := &stripe.CheckoutSessionParams{
		Mode:                stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:          stripe.String(req.SuccessURL),
		CancelURL:           stripe.String(req.CancelURL),
		Customer:            stripe.String(stripeCustomerID),
		ClientReferenceID:   stripe.String(userID.String()),
		AllowPromotionCodes: stripe.Bool(true),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
	}
	checkoutParams.AddMetadata("user_id", userID.String())
	checkoutParams.AddMetadata("plan", string(plan))

	s, err := session.New(checkoutParams)
	if err != nil {
		return nil, fmt.Errorf("create checkout session: %w", err)
	}

	return &dto.CreateCheckoutSessionResponse{
		SessionID:   s.ID,
		CheckoutURL: s.URL,
	}, nil
}

func (uc *subscriptionUseCase) CreatePortalSession(ctx context.Context, userID uuid.UUID, req *dto.CreatePortalSessionRequest) (*dto.CreatePortalSessionResponse, error) {
	if uc.stripeCfg.SecretKey == "" {
		return nil, errors.New("stripe secret key not configured")
	}
	if req == nil || req.ReturnURL == "" {
		return nil, errors.New("return_url is required")
	}

	subscription, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if subscription == nil || subscription.StripeCustomerID == nil || *subscription.StripeCustomerID == "" {
		return nil, errors.New("stripe customer not found for user")
	}

	stripe.Key = uc.stripeCfg.SecretKey

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(*subscription.StripeCustomerID),
		ReturnURL: stripe.String(req.ReturnURL),
	}

	s, err := billingportalsession.New(params)
	if err != nil {
		return nil, fmt.Errorf("create billing portal session: %w", err)
	}

	if s == nil || s.URL == "" {
		return nil, errors.New("billing portal url not returned by stripe")
	}

	return &dto.CreatePortalSessionResponse{PortalURL: s.URL}, nil
}

func (uc *subscriptionUseCase) HandleStripeWebhook(ctx context.Context, payload []byte, signature string) error {
	if uc.stripeCfg.WebhookSecret == "" {
		return errors.New("stripe webhook secret not configured")
	}

	event, err := webhook.ConstructEventWithOptions(payload, signature, uc.stripeCfg.WebhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		return fmt.Errorf("verify stripe webhook signature: %w", err)
	}

	if uc.logger != nil {
		uc.logger.Info(
			"Stripe webhook received",
			zap.String("event_id", event.ID),
			zap.String("event_type", string(event.Type)),
		)
	}

	switch event.Type {
	case "checkout.session.completed":
		return uc.handleCheckoutSessionCompleted(ctx, event)
	case "customer.subscription.created":
		return uc.handleCustomerSubscriptionUpsert(ctx, event)
	case "customer.subscription.updated":
		return uc.handleCustomerSubscriptionUpsert(ctx, event)
	case "invoice.paid":
		return uc.handleInvoicePaid(ctx, event)
	case "invoice.payment_succeeded":
		// Some integrations rely on invoice.payment_succeeded instead of invoice.paid.
		return uc.handleInvoicePaid(ctx, event)
	case "invoice_payment.paid":
		// Observed in local Stripe CLI logs (legacy or compatibility event name).
		return uc.handleInvoicePaid(ctx, event)
	case "invoice.payment_failed":
		return uc.handleInvoicePaymentFailed(ctx, event)
	case "customer.subscription.deleted":
		return uc.handleCustomerSubscriptionDeleted(ctx, event)
	case "charge.refunded":
		return uc.handleChargeRefunded(ctx, event)
	default:
		return nil
	}
}

func (uc *subscriptionUseCase) handleCheckoutSessionCompleted(ctx context.Context, event stripe.Event) error {
	var s stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &s); err != nil {
		return fmt.Errorf("unmarshal checkout.session.completed: %w", err)
	}

	userID, err := uuid.Parse(s.ClientReferenceID)
	if err != nil {
		return fmt.Errorf("invalid client_reference_id (expected user id): %w", err)
	}

	plan := models.SubscriptionPlan(s.Metadata["plan"])
	if !plan.IsValid() || plan == models.SubscriptionPlanFree {
		return fmt.Errorf("invalid plan metadata: %q", s.Metadata["plan"])
	}

	subscription, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if subscription == nil {
		subscription = &models.Subscription{
			UserID: userID,
		}
		if err := uc.subscriptionRepo.Create(ctx, subscription); err != nil {
			return err
		}
	}

	subscription.Plan = plan
	if cid := stripeIDFromCustomer(s.Customer); cid != "" {
		subscription.StripeCustomerID = &cid
	} else {
		subscription.StripeCustomerID = nil
	}
	if sid := stripeIDFromSubscription(s.Subscription); sid != "" {
		subscription.StripeSubscriptionID = &sid
	} else {
		subscription.StripeSubscriptionID = nil
	}
	subscription.Status = models.SubscriptionStatusIncomplete

	if uc.logger != nil {
		uc.logger.Info(
			"Checkout session completed processed",
			zap.String("event_id", event.ID),
			zap.String("user_id", userID.String()),
			zap.String("plan", string(plan)),
			zap.String("stripe_customer_id", strVal(subscription.StripeCustomerID)),
			zap.String("stripe_subscription_id", strVal(subscription.StripeSubscriptionID)),
		)
	}

	return uc.subscriptionRepo.Save(ctx, subscription)
}

func (uc *subscriptionUseCase) handleCustomerSubscriptionUpsert(ctx context.Context, event stripe.Event) error {
	var s stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &s); err != nil {
		return fmt.Errorf("unmarshal customer.subscription: %w", err)
	}
	if s.ID == "" {
		return nil
	}

	stripeCustomerID := stripeIDFromCustomer(s.Customer)

	subscription, err := uc.subscriptionRepo.GetByStripeSubscriptionID(ctx, s.ID)
	if err != nil {
		return err
	}

	if subscription == nil && stripeCustomerID != "" {
		subscription, err = uc.subscriptionRepo.GetByStripeCustomerID(ctx, stripeCustomerID)
		if err != nil {
			return err
		}
	}

	if subscription == nil {
		if uc.logger != nil {
			uc.logger.Warn(
				"Local subscription not found for customer.subscription event, skipping",
				zap.String("event_id", event.ID),
				zap.String("event_type", string(event.Type)),
				zap.String("stripe_subscription_id", s.ID),
				zap.String("stripe_customer_id", stripeCustomerID),
			)
		}
		return nil
	}

	cid := stripeCustomerID
	sid := s.ID
	subscription.StripeCustomerID = &cid
	subscription.StripeSubscriptionID = &sid
	subscription.CancelAtPeriodEnd = s.CancelAtPeriodEnd

	if s.CurrentPeriodEnd > 0 {
		end := time.Unix(s.CurrentPeriodEnd, 0).UTC()
		subscription.CurrentPeriodEnd = &end
	}

	if mapped, ok := mapStripeSubscriptionStatus(s.Status); ok {
		subscription.Status = mapped
	}

	if uc.logger != nil {
		uc.logger.Info(
			"Customer subscription event applied",
			zap.String("event_id", event.ID),
			zap.String("event_type", string(event.Type)),
			zap.String("stripe_subscription_id", s.ID),
			zap.String("stripe_customer_id", stripeCustomerID),
			zap.String("status", string(subscription.Status)),
		)
	}

	return uc.subscriptionRepo.Save(ctx, subscription)
}

func (uc *subscriptionUseCase) handleInvoicePaid(ctx context.Context, event stripe.Event) error {
	var inv stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &inv); err != nil {
		return fmt.Errorf("unmarshal invoice.paid: %w", err)
	}

	stripeSubscriptionID := stripeIDFromSubscription(inv.Subscription)
	stripeCustomerID := stripeIDFromCustomer(inv.Customer)

	var subscription *models.Subscription
	var err error
	if stripeSubscriptionID != "" {
		subscription, err = uc.subscriptionRepo.GetByStripeSubscriptionID(ctx, stripeSubscriptionID)
	} else if stripeCustomerID != "" {
		subscription, err = uc.subscriptionRepo.GetByStripeCustomerID(ctx, stripeCustomerID)
	} else {
		if uc.logger != nil {
			uc.logger.Warn(
				"Stripe invoice event without subscription and customer id, skipping",
				zap.String("event_id", event.ID),
				zap.String("event_type", string(event.Type)),
				zap.String("invoice_id", inv.ID),
			)
		}
		return nil
	}
	if err != nil {
		return err
	}
	if subscription == nil {
		if uc.logger != nil {
			uc.logger.Warn(
				"Local subscription not found for Stripe subscription id, skipping invoice event",
				zap.String("event_id", event.ID),
				zap.String("event_type", string(event.Type)),
				zap.String("stripe_subscription_id", stripeSubscriptionID),
				zap.String("stripe_customer_id", stripeCustomerID),
				zap.String("invoice_id", inv.ID),
			)
		}
		return nil
	}

	now := uc.now()
	subscription.Status = models.SubscriptionStatusActive
	subscription.AccessRevokedAt = nil
	subscription.AccessRevokeReason = nil

	if periodEnd := derivePeriodEndFromInvoice(&inv); periodEnd != nil {
		subscription.CurrentPeriodEnd = periodEnd
	} else if uc.stripeCfg.SecretKey != "" {
		// Fallback: retrieve the subscription from Stripe to get current_period_end reliably.
		stripe.Key = uc.stripeCfg.SecretKey
		if stripeSubscriptionID != "" {
			if stripeSub, err := stripesub.Get(stripeSubscriptionID, nil); err == nil && stripeSub != nil {
				if stripeSub.CurrentPeriodEnd > 0 {
					end := time.Unix(stripeSub.CurrentPeriodEnd, 0).UTC()
					subscription.CurrentPeriodEnd = &end
				}
				subscription.CancelAtPeriodEnd = stripeSub.CancelAtPeriodEnd
				if mapped, ok := mapStripeSubscriptionStatus(stripeSub.Status); ok {
					subscription.Status = mapped
				}
			} else if uc.logger != nil && err != nil {
				uc.logger.Warn(
					"Failed to fetch Stripe subscription for period end fallback",
					zap.String("stripe_subscription_id", stripeSubscriptionID),
					zap.String("invoice_id", inv.ID),
					zap.Error(err),
				)
			}
		} else if uc.logger != nil {
			uc.logger.Warn(
				"Stripe subscription id missing; cannot fetch period end fallback from Stripe subscription",
				zap.String("stripe_customer_id", stripeCustomerID),
				zap.String("invoice_id", inv.ID),
			)
		}
	}

	// If Stripe provided cancellation info, keep it updated.
	if inv.Subscription != nil && inv.Subscription.CancelAtPeriodEnd {
		subscription.CancelAtPeriodEnd = true
	} else if subscription.CancelAtPeriodEnd {
		subscription.CancelAtPeriodEnd = false
	}

	// Clear canceled time on renewal if applicable.
	if subscription.CanceledAt != nil && now.Before(*subscription.CanceledAt) {
		subscription.CanceledAt = nil
	}

	return uc.subscriptionRepo.Save(ctx, subscription)
}

func derivePeriodEndFromInvoice(inv *stripe.Invoice) *time.Time {
	if inv == nil {
		return nil
	}

	// Prefer invoice line item periods (most accurate for subscriptions).
	periodEnd := time.Time{}
	if inv.Lines != nil {
		for _, line := range inv.Lines.Data {
			if line.Period == nil || line.Period.End <= 0 {
				continue
			}

			end := time.Unix(line.Period.End, 0).UTC()
			if end.After(periodEnd) {
				periodEnd = end
			}
		}
	}
	if !periodEnd.IsZero() {
		return &periodEnd
	}

	// Fallback: use invoice period end if present.
	if inv.PeriodEnd > 0 {
		end := time.Unix(inv.PeriodEnd, 0).UTC()
		return &end
	}

	return nil
}

func strVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func stripeIDFromCustomer(c *stripe.Customer) string {
	if c == nil {
		return ""
	}
	return c.ID
}

func stripeIDFromSubscription(s *stripe.Subscription) string {
	if s == nil {
		return ""
	}
	return s.ID
}

func mapStripeSubscriptionStatus(status stripe.SubscriptionStatus) (models.SubscriptionStatus, bool) {
	switch status {
	case stripe.SubscriptionStatusActive:
		return models.SubscriptionStatusActive, true
	case stripe.SubscriptionStatusTrialing:
		return models.SubscriptionStatusTrialing, true
	case stripe.SubscriptionStatusPastDue:
		return models.SubscriptionStatusPastDue, true
	case stripe.SubscriptionStatusCanceled:
		return models.SubscriptionStatusCanceled, true
	case stripe.SubscriptionStatusIncomplete:
		return models.SubscriptionStatusIncomplete, true
	case stripe.SubscriptionStatusIncompleteExpired:
		return models.SubscriptionStatusIncompleteExpired, true
	case stripe.SubscriptionStatusUnpaid:
		return models.SubscriptionStatusUnpaid, true
	case stripe.SubscriptionStatusPaused:
		return models.SubscriptionStatusPaused, true
	default:
		return "", false
	}
}

func (uc *subscriptionUseCase) handleInvoicePaymentFailed(ctx context.Context, event stripe.Event) error {
	var inv stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &inv); err != nil {
		return fmt.Errorf("unmarshal invoice.payment_failed: %w", err)
	}
	if inv.Subscription == nil || inv.Subscription.ID == "" {
		return nil
	}

	subscription, err := uc.subscriptionRepo.GetByStripeSubscriptionID(ctx, inv.Subscription.ID)
	if err != nil {
		return err
	}
	if subscription == nil {
		return nil
	}

	subscription.Status = models.SubscriptionStatusPastDue
	return uc.subscriptionRepo.Save(ctx, subscription)
}

func (uc *subscriptionUseCase) handleCustomerSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var s stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &s); err != nil {
		return fmt.Errorf("unmarshal customer.subscription.deleted: %w", err)
	}
	if s.ID == "" {
		return nil
	}

	subscription, err := uc.subscriptionRepo.GetByStripeSubscriptionID(ctx, s.ID)
	if err != nil {
		return err
	}
	if subscription == nil {
		return nil
	}

	now := uc.now()
	subscription.Status = models.SubscriptionStatusCanceled
	subscription.Plan = models.SubscriptionPlanFree
	subscription.CanceledAt = &now
	subscription.CancelAtPeriodEnd = false
	subscription.CurrentPeriodEnd = nil

	return uc.subscriptionRepo.Save(ctx, subscription)
}

func (uc *subscriptionUseCase) handleChargeRefunded(ctx context.Context, event stripe.Event) error {
	var c stripe.Charge
	if err := json.Unmarshal(event.Data.Raw, &c); err != nil {
		return fmt.Errorf("unmarshal charge.refunded: %w", err)
	}
	if c.Customer == nil || c.Customer.ID == "" {
		return nil
	}

	subscription, err := uc.subscriptionRepo.GetByStripeCustomerID(ctx, c.Customer.ID)
	if err != nil {
		return err
	}
	if subscription == nil {
		return nil
	}

	now := uc.now()
	reason := "refunded"
	subscription.AccessRevokedAt = &now
	subscription.AccessRevokeReason = &reason
	subscription.Status = models.SubscriptionStatusAccessRevokedRefund

	return uc.subscriptionRepo.Save(ctx, subscription)
}

func (uc *subscriptionUseCase) priceIDForPlan(plan models.SubscriptionPlan) (string, error) {
	switch plan {
	case models.SubscriptionPlanSimple:
		if uc.stripeCfg.PriceSimple == "" {
			return "", errors.New("stripe price id for simple plan not configured")
		}
		return uc.stripeCfg.PriceSimple, nil
	case models.SubscriptionPlanMedium:
		if uc.stripeCfg.PriceMedium == "" {
			return "", errors.New("stripe price id for medium plan not configured")
		}
		return uc.stripeCfg.PriceMedium, nil
	case models.SubscriptionPlanUltra:
		if uc.stripeCfg.PriceUltra == "" {
			return "", errors.New("stripe price id for ultra plan not configured")
		}
		return uc.stripeCfg.PriceUltra, nil
	default:
		return "", fmt.Errorf("unsupported plan: %s", plan)
	}
}

func isSubscriptionActive(now time.Time, subscription *models.Subscription) (bool, string) {
	if subscription.AccessRevokedAt != nil {
		return false, "access_revoked"
	}

	switch subscription.Status {
	case models.SubscriptionStatusActive, models.SubscriptionStatusTrialing:
	default:
		return false, "stripe_status_not_active"
	}

	if subscription.CurrentPeriodEnd != nil && now.After(*subscription.CurrentPeriodEnd) {
		return false, "current_period_ended"
	}

	return true, "ok"
}
