package paypal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gitee.com/meinongyihe/travel-api/internal/paypal"
	"gitee.com/meinongyihe/travel-api/internal/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

// Adapter wraps the legacy paypal.Provider and implements payment.PaymentProvider.
type Adapter struct {
	inner *paypal.Provider
	cfg   paypal.Config
}

// New creates a PayPal PaymentProvider adapter.
func New(cfg paypal.Config) *Adapter {
	return &Adapter{inner: paypal.New(cfg), cfg: cfg}
}

// Inner returns the underlying legacy paypal.Provider for backward-compatible callers.
func (a *Adapter) Inner() *paypal.Provider { return a.inner }

func (a *Adapter) Name() string { return "paypal" }

// CreateCheckout delegates to the legacy CreateOrder.
func (a *Adapter) CreateCheckout(ctx context.Context, in payment.CheckoutInput) (*payment.CheckoutResult, error) {
	checkout, err := a.inner.CreateOrder(ctx, paypal.CheckoutInput{
		PaymentNo: in.PaymentNo,
		Amount:    in.Amount,
		Currency:  in.Currency,
		State:     in.State,
	})
	if err != nil {
		return nil, err
	}
	return &payment.CheckoutResult{
		CheckoutURL:       checkout.CheckoutURL,
		ProviderPaymentID: checkout.OrderID,
	}, nil
}

// ConfirmPayment delegates to the legacy CaptureOrder.
func (a *Adapter) ConfirmPayment(ctx context.Context, providerPaymentID string) (*payment.ConfirmResult, error) {
	cap, err := a.inner.CaptureOrder(ctx, providerPaymentID)
	if err != nil {
		return nil, err
	}
	status := "PENDING"
	switch strings.ToUpper(cap.Status) {
	case "COMPLETED", "CAPTURED":
		status = "PAID"
	case "FAILED", "DECLINED":
		status = "FAILED"
	}
	return &payment.ConfirmResult{
		ProviderPaymentID: cap.OrderID,
		CaptureID:         cap.CaptureID,
		Status:            status,
	}, nil
}

// Refund delegates to the legacy RefundCapture.
func (a *Adapter) Refund(ctx context.Context, providerPaymentID string, amount int64, currency string) (*payment.RefundResult, error) {
	r, err := a.inner.RefundCapture(ctx, providerPaymentID, amount, currency)
	if err != nil {
		return nil, err
	}
	status := "PENDING"
	switch strings.ToUpper(r.Status) {
	case "COMPLETED", "SUCCEEDED":
		status = "SUCCEEDED"
	case "FAILED":
		status = "FAILED"
	}
	return &payment.RefundResult{RefundID: r.RefundID, Status: status}, nil
}

// VerifyWebhook validates the PayPal webhook and returns a normalised event.
func (a *Adapter) VerifyWebhook(ctx context.Context, headers http.Header, body []byte) (payment.WebhookEvent, error) {
	if err := a.inner.VerifyWebhook(ctx, headers, json.RawMessage(body)); err != nil {
		return payment.WebhookEvent{}, err
	}
	// Parse the PayPal event envelope.
	var envelope struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID             string `json:"id"`
			CustomID       string `json:"custom_id"`
			InvoiceID      string `json:"invoice_id"`
			Amount         struct {
				Value        string `json:"value"`
				CurrencyCode string `json:"currency_code"`
			} `json:"amount"`
			Status string `json:"status"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return payment.WebhookEvent{}, fmt.Errorf("paypal webhook unmarshal: %w", err)
	}

	eventType := payment.EventPaymentFailed
	switch envelope.EventType {
	case "PAYMENT.CAPTURE.COMPLETED":
		eventType = payment.EventPaymentSucceeded
	case "PAYMENT.CAPTURE.DENIED", "PAYMENT.CAPTURE.REFUNDED":
		eventType = payment.EventPaymentFailed
	}

	// Extract payment_no from custom_id (state token) or invoice_id.
	paymentNo := envelope.Resource.InvoiceID
	if paymentNo == "" && envelope.Resource.CustomID != "" {
		// custom_id contains the signed state; parse it to get payment_no.
		if _, _, pn, err := paypal.ParseState(a.cfg.Secret, envelope.Resource.CustomID); err == nil {
			paymentNo = pn
		}
	}

	logx.WithContext(ctx).Infof("[PayPal Webhook] event=%s payment_no=%s id=%s",
		envelope.EventType, paymentNo, envelope.Resource.ID)

	return payment.WebhookEvent{
		Type:              eventType,
		ProviderPaymentID: envelope.Resource.ID,
		PaymentNo:         paymentNo,
		Raw:               json.RawMessage(body),
	}, nil
}
