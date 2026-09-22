package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gitee.com/meinongyihe/travel-api/internal/payment"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/refund"
	"github.com/stripe/stripe-go/v82/webhook"
	"github.com/zeromicro/go-zero/core/logx"
)

// Config holds Stripe API credentials.
type Config struct {
	Enabled        bool
	SecretKey      string // sk_test_... or sk_live_...
	PublishableKey string // pk_test_... or pk_live_...
	WebhookSecret  string // whsec_...
}

// Provider implements payment.PaymentProvider for Stripe.
type Provider struct {
	cfg Config
}

// New creates a Stripe PaymentProvider.
func New(cfg Config) *Provider {
	// Set the global Stripe key so that package-level functions work.
	if cfg.SecretKey != "" {
		stripe.Key = cfg.SecretKey
	}
	return &Provider{cfg: cfg}
}

func (p *Provider) Name() string { return "stripe" }

// PublishableKey returns the publishable key for frontend use.
func (p *Provider) PublishableKey() string { return p.cfg.PublishableKey }

// CreateCheckout creates a Stripe PaymentIntent and returns the client_secret
// so the frontend can confirm payment via Stripe Elements.
func (p *Provider) CreateCheckout(ctx context.Context, in payment.CheckoutInput) (*payment.CheckoutResult, error) {
	if in.Amount <= 0 {
		return nil, fmt.Errorf("stripe: invalid amount %d", in.Amount)
	}
	if in.Currency == "" {
		return nil, fmt.Errorf("stripe: currency is required")
	}

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(in.Amount),
		Currency:           stripe.String(strings.ToLower(in.Currency)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Description:        stripe.String(fmt.Sprintf("Travel Order %s", in.OrderNo)),
		Metadata: map[string]string{
			"payment_no": in.PaymentNo,
			"order_no":   in.OrderNo,
		},
	}
	// Capture manually so we can control when funds are collected.
	params.CaptureMethod = stripe.String("automatic")

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe create payment intent: %w", err)
	}

	return &payment.CheckoutResult{
		ClientSecret:      pi.ClientSecret,
		ProviderPaymentID: pi.ID,
	}, nil
}

// ConfirmPayment retrieves a PaymentIntent to check its status.
// With Stripe Elements the frontend confirms payment via stripe.confirmPayment();
// this method is used server-side to verify the final status.
func (p *Provider) ConfirmPayment(ctx context.Context, providerPaymentID string) (*payment.ConfirmResult, error) {
	pi, err := paymentintent.Get(providerPaymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe get payment intent: %w", err)
	}

	status := "PENDING"
	switch pi.Status {
	case stripe.PaymentIntentStatusSucceeded:
		status = "PAID"
	case stripe.PaymentIntentStatusCanceled:
		status = "FAILED"
	case stripe.PaymentIntentStatusRequiresPaymentMethod,
		stripe.PaymentIntentStatusRequiresConfirmation:
		status = "PENDING"
	case stripe.PaymentIntentStatusProcessing,
		stripe.PaymentIntentStatusRequiresAction:
		status = "PENDING"
	}

	return &payment.ConfirmResult{
		ProviderPaymentID: pi.ID,
		Status:            status,
	}, nil
}

// Refund creates a full or partial refund against a PaymentIntent.
func (p *Provider) Refund(ctx context.Context, providerPaymentID string, amount int64, currency string) (*payment.RefundResult, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(providerPaymentID),
	}
	if amount > 0 {
		params.Amount = stripe.Int64(amount)
	}

	r, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe refund: %w", err)
	}

	status := "PENDING"
	switch r.Status {
	case stripe.RefundStatusSucceeded:
		status = "SUCCEEDED"
	case stripe.RefundStatusFailed:
		status = "FAILED"
	case stripe.RefundStatusPending:
		status = "PENDING"
	}

	return &payment.RefundResult{RefundID: r.ID, Status: status}, nil
}

// VerifyWebhook validates the Stripe webhook signature and returns a normalised event.
func (p *Provider) VerifyWebhook(ctx context.Context, headers http.Header, body []byte) (payment.WebhookEvent, error) {
	sigHeader := headers.Get("Stripe-Signature")
	if sigHeader == "" {
		return payment.WebhookEvent{}, fmt.Errorf("stripe: missing Stripe-Signature header")
	}

	event, err := webhook.ConstructEvent(body, sigHeader, p.cfg.WebhookSecret)
	if err != nil {
		return payment.WebhookEvent{}, fmt.Errorf("stripe webhook signature verification failed: %w", err)
	}

	logx.WithContext(ctx).Infof("[Stripe Webhook] type=%s id=%s", event.Type, event.ID)

	evt := payment.WebhookEvent{
		Raw: json.RawMessage(body),
	}

	switch event.Type {
	case "payment_intent.succeeded":
		evt.Type = payment.EventPaymentSucceeded
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return evt, fmt.Errorf("stripe unmarshal payment_intent: %w", err)
		}
		evt.ProviderPaymentID = pi.ID
		evt.Amount = pi.Amount
		evt.Currency = strings.ToUpper(string(pi.Currency))
		if pn, ok := pi.Metadata["payment_no"]; ok {
			evt.PaymentNo = pn
		}

	case "payment_intent.payment_failed":
		evt.Type = payment.EventPaymentFailed
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return evt, fmt.Errorf("stripe_unmarshal failed_pi: %w", err)
		}
		evt.ProviderPaymentID = pi.ID
		if pn, ok := pi.Metadata["payment_no"]; ok {
			evt.PaymentNo = pn
		}

	default:
		// Unhandled event type — return as-is so caller can decide.
		evt.Type = string(event.Type)
	}

	return evt, nil
}
