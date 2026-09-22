package shopify

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gitee.com/meinongyihe/travel-api/internal/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

// Config holds Shopify API credentials and store information.
type Config struct {
	Enabled       bool
	StoreDomain   string // e.g. "your-store.myshopify.com"
	AccessToken   string // Admin API access token (shpat_...)
	APIVersion    string // e.g. "2024-01"
	WebhookSecret string // HMAC shared secret for webhook verification
	SyncEnabled   bool   // whether product sync is enabled
}

// Provider implements payment.PaymentProvider for Shopify as a sales channel.
type Provider struct {
	cfg        Config
	httpClient *http.Client
}

// New creates a Shopify PaymentProvider.
func New(cfg Config) *Provider {
	if cfg.APIVersion == "" {
		cfg.APIVersion = "2024-01"
	}
	return &Provider{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *Provider) Name() string { return "shopify" }

// CreateCheckout is a no-op for Shopify — checkout happens on the Shopify storefront.
// This method exists to satisfy the PaymentProvider interface.
func (p *Provider) CreateCheckout(_ context.Context, _ payment.CheckoutInput) (*payment.CheckoutResult, error) {
	return nil, fmt.Errorf("shopify: checkout is handled by the Shopify storefront; use product sync instead")
}

// ConfirmPayment is a no-op for Shopify — payment confirmation comes via webhook.
func (p *Provider) ConfirmPayment(_ context.Context, _ string) (*payment.ConfirmResult, error) {
	return nil, fmt.Errorf("shopify: payment confirmation is handled via Shopify webhook")
}

// Refund creates a refund via Shopify Admin API against an order.
func (p *Provider) Refund(ctx context.Context, providerPaymentID string, amount int64, currency string) (*payment.RefundResult, error) {
	// providerPaymentID for Shopify is the Shopify Order ID.
	if providerPaymentID == "" {
		return nil, fmt.Errorf("shopify: order ID is required for refund")
	}

	// First, get the transaction ID for the order.
	txURL := p.adminURL(fmt.Sprintf("/orders/%s/transactions.json", providerPaymentID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, txURL, nil)
	if err != nil {
		return nil, err
	}
	p.setHeaders(req)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("shopify get transactions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("shopify get transactions failed: status=%d body=%s", resp.StatusCode, truncate(string(body)))
	}

	var txResp struct {
		Transactions []struct {
			ID     int64  `json:"id"`
			Kind   string `json:"kind"`
			Status string `json:"status"`
		} `json:"transactions"`
	}
	data, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &txResp); err != nil {
		return nil, err
	}

	// Find the capture/sale transaction to refund against.
	var transactionID int64
	for _, tx := range txResp.Transactions {
		if tx.Kind == "capture" || tx.Kind == "sale" {
			transactionID = tx.ID
			break
		}
	}
	if transactionID == 0 {
		return nil, fmt.Errorf("shopify: no capture/sale transaction found for order %s", providerPaymentID)
	}

	// Create the refund.
	refundBody := map[string]any{
		"refund": map[string]any{
			"notify":      true,
			"transactions": []map[string]any{
				{
					"parent_id": transactionID,
					"kind":      "refund",
					"amount":    minorToMajor(amount, currency),
				},
			},
		},
	}
	raw, _ := json.Marshal(refundBody)
	refundURL := p.adminURL(fmt.Sprintf("/orders/%s/refunds.json", providerPaymentID))
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, refundURL, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	p.setHeaders(req)

	resp, err = p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("shopify create refund: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("shopify create refund failed: status=%d body=%s", resp.StatusCode, truncate(string(body)))
	}

	var refundResp struct {
		ID     int64  `json:"id"`
		Status string `json:"status"`
	}
	respData, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respData, &refundResp); err != nil {
		return nil, err
	}

	return &payment.RefundResult{
		RefundID: fmt.Sprintf("%d", refundResp.ID),
		Status:   "SUCCEEDED",
	}, nil
}

// VerifyWebhook validates the Shopify webhook HMAC-SHA256 signature and returns a normalised event.
func (p *Provider) VerifyWebhook(ctx context.Context, headers http.Header, body []byte) (payment.WebhookEvent, error) {
	evt := payment.WebhookEvent{Raw: json.RawMessage(body)}

	// Verify HMAC-SHA256 signature.
	mac := hmac.New(sha256.New, []byte(p.cfg.WebhookSecret))
	mac.Write(body)
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	receivedSig := headers.Get("X-Shopify-Hmac-Sha256")
	if receivedSig == "" {
		return evt, fmt.Errorf("shopify: missing X-Shopify-Hmac-Sha256 header")
	}
	if !hmac.Equal([]byte(expectedSig), []byte(receivedSig)) {
		return evt, fmt.Errorf("shopify: webhook signature mismatch")
	}

	// Parse the topic from the header.
	topic := headers.Get("X-Shopify-Topic")
	logx.WithContext(ctx).Infof("[Shopify Webhook] topic=%s", topic)

	switch topic {
	case "orders/create":
		evt.Type = payment.EventOrderCreated
		var order struct {
			ID               int64  `json:"id"`
			OrderNumber      int64  `json:"order_number"`
			FinancialStatus  string `json:"financial_status"`
			TotalPrice       string `json:"total_price"`
			Currency         string `json:"currency"`
			Note             string `json:"note"`
			NoteAttributes   []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"note_attributes"`
		}
		if err := json.Unmarshal(body, &order); err != nil {
			return evt, fmt.Errorf("shopify unmarshal order: %w", err)
		}
		evt.ProviderPaymentID = fmt.Sprintf("%d", order.ID)
		// Extract payment_no from note_attributes if present.
		for _, attr := range order.NoteAttributes {
			if strings.EqualFold(attr.Name, "payment_no") {
				evt.PaymentNo = attr.Value
				break
			}
		}

	case "orders/paid":
		evt.Type = payment.EventPaymentSucceeded
		var order struct {
			ID         int64  `json:"id"`
			TotalPrice string `json:"total_price"`
			Currency   string `json:"currency"`
			NoteAttributes []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"note_attributes"`
		}
		if err := json.Unmarshal(body, &order); err != nil {
			return evt, fmt.Errorf("shopify unmarshal paid order: %w", err)
		}
		evt.ProviderPaymentID = fmt.Sprintf("%d", order.ID)
		for _, attr := range order.NoteAttributes {
			if strings.EqualFold(attr.Name, "payment_no") {
				evt.PaymentNo = attr.Value
				break
			}
		}

	case "orders/cancelled":
		evt.Type = payment.EventOrderCancelled
		var order struct {
			ID int64 `json:"id"`
		}
		if err := json.Unmarshal(body, &order); err != nil {
			return evt, fmt.Errorf("shopify unmarshal cancelled order: %w", err)
		}
		evt.ProviderPaymentID = fmt.Sprintf("%d", order.ID)

	default:
		evt.Type = topic
	}

	return evt, nil
}

// ── Product Sync (Shopify-specific, not part of PaymentProvider interface) ──

// ProductInput is the data needed to sync a travel product to Shopify.
type ProductInput struct {
	Title       string
	Description string
	Vendor      string
	ProductType string
	Tags        []string
	Images      []string // image URLs
	Variants    []VariantInput
}

// VariantInput represents a product variant (package/SKU).
type VariantInput struct {
	Title    string
	Price    string // major units as string, e.g. "99.00"
	SKU      string
	Quantity int
}

// SyncProduct creates or updates a Shopify product from travel product data.
func (p *Provider) SyncProduct(ctx context.Context, in ProductInput) (shopifyProductID string, err error) {
	if !p.cfg.SyncEnabled {
		return "", fmt.Errorf("shopify: product sync is disabled")
	}

	// Build Shopify product payload.
	variants := make([]map[string]any, 0, len(in.Variants))
	for _, v := range in.Variants {
		variants = append(variants, map[string]any{
			"title":     v.Title,
			"price":     v.Price,
			"sku":       v.SKU,
			"inventory_quantity": v.Quantity,
			"inventory_management": "shopify",
		})
	}

	images := make([]map[string]any, 0, len(in.Images))
	for _, img := range in.Images {
		images = append(images, map[string]any{"src": img})
	}

	payload := map[string]any{
		"product": map[string]any{
			"title":        in.Title,
			"body_html":    in.Description,
			"vendor":       in.Vendor,
			"product_type": in.ProductType,
			"tags":         strings.Join(in.Tags, ", "),
			"variants":     variants,
			"images":       images,
		},
	}

	raw, _ := json.Marshal(payload)
	url := p.adminURL("/products.json")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	p.setHeaders(req)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("shopify sync product: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("shopify sync product failed: status=%d body=%s", resp.StatusCode, truncate(string(body)))
	}

	var result struct {
		Product struct {
			ID int64 `json:"id"`
		} `json:"product"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", result.Product.ID), nil
}

// ── Internal helpers ──

func (p *Provider) adminURL(path string) string {
	return fmt.Sprintf("https://%s/admin/api/%s%s", p.cfg.StoreDomain, p.cfg.APIVersion, path)
}

func (p *Provider) setHeaders(req *http.Request) {
	req.Header.Set("X-Shopify-Access-Token", p.cfg.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

func minorToMajor(amount int64, currency string) string {
	switch strings.ToUpper(currency) {
	case "JPY", "KRW":
		return fmt.Sprintf("%d", amount)
	default:
		return fmt.Sprintf("%.2f", float64(amount)/100)
	}
}

func truncate(s string) string {
	if len(s) > 500 {
		return s[:500] + "..."
	}
	return s
}
