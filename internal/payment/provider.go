package payment

import (
	"context"
	"encoding/json"
	"net/http"
)

// PaymentProvider is the unified interface for all payment channels (PayPal / Stripe / Shopify).
// Each provider implements checkout creation, payment confirmation, refund, and webhook verification.
type PaymentProvider interface {
	// Name returns the canonical provider identifier ("paypal", "stripe", "shopify").
	Name() string

	// CreateCheckout initiates a payment and returns whatever the caller needs to complete it.
	// For redirect-based providers (PayPal, Shopify) this returns a CheckoutURL.
	// For embedded providers (Stripe Elements) this returns a ClientSecret.
	CreateCheckout(ctx context.Context, in CheckoutInput) (*CheckoutResult, error)

	// ConfirmPayment explicitly captures / confirms a payment that is in a pending state.
	// Used by redirect-flow providers where the user returns to our site after paying.
	ConfirmPayment(ctx context.Context, providerPaymentID string) (*ConfirmResult, error)

	// Refund issues a full or partial refund against a previously captured payment.
	Refund(ctx context.Context, providerPaymentID string, amount int64, currency string) (*RefundResult, error)

	// VerifyWebhook validates an incoming webhook request and extracts a normalised event.
	VerifyWebhook(ctx context.Context, headers http.Header, body []byte) (WebhookEvent, error)
}

// CheckoutInput carries everything a provider needs to start a payment.
type CheckoutInput struct {
	PaymentNo  string
	OrderNo    string
	Amount     int64
	Currency   string
	SuccessURL string // redirect URL after successful payment (PayPal / Shopify)
	CancelURL  string // redirect URL if user cancels (PayPal)
	State      string // signed state token for callback verification (PayPal)
	Metadata   map[string]string
}

// CheckoutResult is returned by CreateCheckout.
// Redirect-based providers fill CheckoutURL; embedded providers fill ClientSecret.
type CheckoutResult struct {
	CheckoutURL       string // PayPal / Shopify redirect target
	ClientSecret      string // Stripe PaymentIntent client_secret
	ProviderPaymentID string // provider-side payment identifier
}

// ConfirmResult is returned by ConfirmPayment (capture).
type ConfirmResult struct {
	ProviderPaymentID string
	CaptureID         string // provider-specific capture / refund reference
	Status            string // normalised status: "PAID", "FAILED", "PENDING"
}

// RefundResult is returned by Refund.
type RefundResult struct {
	RefundID string
	Status   string // "SUCCEEDED", "PENDING", "FAILED"
}

// WebhookEvent is a normalised webhook payload shared across providers.
type WebhookEvent struct {
	Type              string          // "payment.succeeded" | "payment.failed" | "order.created" | ...
	ProviderPaymentID string          // provider-side payment ID
	PaymentNo         string          // our payment_no (extracted from metadata / custom_id / state)
	Amount            int64           // amount in minor units
	Currency          string          // ISO-4217 currency code
	Raw               json.RawMessage // original webhook payload for provider-specific processing
}

// Webhook event type constants — normalised across providers.
const (
	EventPaymentSucceeded = "payment.succeeded"
	EventPaymentFailed    = "payment.failed"
	EventOrderCreated     = "order.created"
	EventOrderPaid        = "order.paid"
	EventOrderCancelled   = "order.cancelled"
)
