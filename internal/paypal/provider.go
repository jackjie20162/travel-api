package paypal

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type Config struct {
	Enabled                                                    bool
	BaseURL, ClientID, Secret, WebhookID, ReturnURL, CancelURL string
	SkipVerify                                                 bool
}
type Provider struct {
	cfg        Config
	httpClient *http.Client
}
type CheckoutInput struct {
	PaymentNo string
	Amount    int64
	Currency  string
	State     string
}
type Checkout struct {
	OrderID     string
	CheckoutURL string
}
type Capture struct {
	OrderID   string
	CaptureID string
	Status    string
}
type RefundResult struct {
	RefundID string
	Status   string
}

func New(cfg Config) *Provider {
	if cfg.BaseURL == "" {
		if strings.EqualFold(os.Getenv("PAYPAL_ENV"), "live") {
			cfg.BaseURL = "https://api-m.paypal.com"
		} else {
			cfg.BaseURL = "https://api-m.sandbox.paypal.com"
		}
	}
	if cfg.ClientID == "" {
		cfg.ClientID = os.Getenv("PAYPAL_CLIENT_ID")
	}
	if cfg.Secret == "" {
		cfg.Secret = os.Getenv("PAYPAL_CLIENT_SECRET")
	}
	if cfg.WebhookID == "" {
		cfg.WebhookID = os.Getenv("PAYPAL_WEBHOOK_ID")
	}
	return &Provider{cfg: cfg, httpClient: &http.Client{Timeout: 15 * time.Second}}
}

func (p *Provider) CreateOrder(ctx context.Context, in CheckoutInput) (*Checkout, error) {
	if p.cfg.ClientID == "" || p.cfg.Secret == "" {
		return nil, errors.New("paypal credentials are not configured")
	}
	if in.Amount <= 0 || in.Currency == "" || in.PaymentNo == "" || in.State == "" {
		return nil, errors.New("invalid paypal checkout input")
	}
	token, err := p.accessToken(ctx)
	if err != nil {
		return nil, err
	}
	returnURL := p.cfg.ReturnURL
	if returnURL == "" {
		returnURL = "https://travel.code688.com/api/travel/payments/paypal/return"
	}
	cancelURL := p.cfg.CancelURL
	if cancelURL == "" {
		cancelURL = "https://travel.code688.com/api/travel/payments/paypal/cancel"
	}
	returnURL = addQuery(returnURL, "state", in.State)
	cancelURL = addQuery(cancelURL, "state", in.State)
	body := map[string]any{"intent": "CAPTURE", "purchase_units": []any{map[string]any{"reference_id": in.PaymentNo, "custom_id": in.State, "invoice_id": in.PaymentNo, "amount": map[string]string{"currency_code": strings.ToUpper(in.Currency), "value": minorToMajor(in.Amount, in.Currency)}}}, "application_context": map[string]any{"return_url": returnURL, "cancel_url": cancelURL, "user_action": "PAY_NOW"}}
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(p.cfg.BaseURL, "/")+"/v2/checkout/orders", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PayPal-Request-Id", in.PaymentNo)
	req.Header.Set("Prefer", "return=representation")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("paypal create order failed: status=%d body=%s", resp.StatusCode, truncate(string(data)))
	}
	var out struct {
		ID    string `json:"id"`
		Links []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, errors.New("paypal create order returned empty id")
	}
	for _, link := range out.Links {
		if link.Rel == "approve" || link.Rel == "payer-action" {
			return &Checkout{OrderID: out.ID, CheckoutURL: link.Href}, nil
		}
	}
	return nil, errors.New("paypal create order returned no approval url")
}

func (p *Provider) CaptureOrder(ctx context.Context, orderID string) (*Capture, error) {
	token, err := p.accessToken(ctx)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(p.cfg.BaseURL, "/")+"/v2/checkout/orders/"+url.PathEscape(orderID)+"/capture", bytes.NewReader([]byte(`{}`)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("paypal capture failed: status=%d body=%s", resp.StatusCode, truncate(string(data)))
	}
	var out struct {
		ID            string `json:"id"`
		Status        string `json:"status"`
		PurchaseUnits []struct {
			Payments struct {
				Captures []struct {
					ID     string `json:"id"`
					Status string `json:"status"`
				} `json:"captures"`
			} `json:"payments"`
		} `json:"purchase_units"`
	}

	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	c := &Capture{OrderID: out.ID, Status: out.Status}
	if len(out.PurchaseUnits) > 0 && len(out.PurchaseUnits[0].Payments.Captures) > 0 {
		c.CaptureID = out.PurchaseUnits[0].Payments.Captures[0].ID
		if c.Status == "" {
			c.Status = out.PurchaseUnits[0].Payments.Captures[0].Status
		}
	}
	return c, nil
}

func (p *Provider) RefundCapture(ctx context.Context, captureID string, amount int64, currency string) (*RefundResult, error) {
	token, err := p.accessToken(ctx)
	if err != nil {
		return nil, err
	}
	body := map[string]any{}
	if amount > 0 {
		body["amount"] = map[string]string{"value": minorToMajor(amount, currency), "currency_code": strings.ToUpper(currency)}
	}
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(p.cfg.BaseURL, "/")+"/v2/payments/captures/"+url.PathEscape(captureID)+"/refund", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("paypal refund failed: status=%d body=%s", resp.StatusCode, truncate(string(data)))
	}
	var out struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &RefundResult{RefundID: out.ID, Status: out.Status}, nil
}

func (p *Provider) VerifyWebhook(ctx context.Context, headers http.Header, event json.RawMessage) error {
	if p.cfg.SkipVerify {
		logx.Infof("[PayPal Webhook] signature verification skipped (SkipVerify=true)")
		return nil
	}
	if p.cfg.WebhookID == "" {
		return errors.New("paypal webhook id is not configured")
	}
	// Use PayPal's official verify-webhook-signature API
	return p.verifyWebhookAPI(ctx, headers, event)
}

// verifyWebhookAPI calls PayPal's verify-webhook-signature API.
func (p *Provider) verifyWebhookAPI(ctx context.Context, headers http.Header, event json.RawMessage) error {
	authAlgo := headers.Get("PAYPAL-AUTH-ALGO")
	certURL := headers.Get("PAYPAL-CERT-URL")
	transmissionID := headers.Get("PAYPAL-TRANSMISSION-ID")
	transmissionSig := headers.Get("PAYPAL-TRANSMISSION-SIG")
	transmissionTime := headers.Get("PAYPAL-TRANSMISSION-TIME")
	token, err := p.accessToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}
	payload := map[string]any{"auth_algo": authAlgo, "cert_url": certURL, "transmission_id": transmissionID, "transmission_sig": transmissionSig, "transmission_time": transmissionTime, "webhook_id": p.cfg.WebhookID}
	var obj any
	if err := json.Unmarshal(event, &obj); err != nil {
		return err
	}
	payload["webhook_event"] = obj
	raw, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(p.cfg.BaseURL, "/")+"/v1/notifications/verify-webhook-signature", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("verification API error: status=%d body=%s", resp.StatusCode, truncate(string(data)))
	}
	var out struct {
		VerificationStatus string `json:"verification_status"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return fmt.Errorf("failed to parse verification response: %w, body=%s", err, truncate(string(data)))
	}
	if out.VerificationStatus != "SUCCESS" {
		return fmt.Errorf("verification failed: status=%s webhook_id=%q",
			out.VerificationStatus, p.cfg.WebhookID)
	}
	return nil
}

func (p *Provider) accessToken(ctx context.Context) (string, error) {
	if p.cfg.ClientID == "" || p.cfg.Secret == "" {
		return "", errors.New("paypal credentials are not configured")
	}
	form := url.Values{"grant_type": {"client_credentials"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(p.cfg.BaseURL, "/")+"/v1/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(p.cfg.ClientID, p.cfg.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("paypal oauth failed: status=%d body=%s", resp.StatusCode, truncate(string(data)))
	}
	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return "", err
	}
	if out.AccessToken == "" {
		return "", errors.New("paypal oauth returned empty access token")
	}
	return out.AccessToken, nil
}
func minorToMajor(amount int64, currency string) string {
	switch strings.ToUpper(currency) {
	case "JPY", "KRW":
		return strconv.FormatInt(amount, 10)
	default:
		// 后端金额单位为元（int64），直接格式化为两位小数
		return fmt.Sprintf("%d.00", amount)
	}
}
func addQuery(rawURL, key, value string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	return u.String()
}
func truncate(s string) string {
	return truncateS(s, 1000)
}
func truncateS(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func SignState(secret string, tenantID, merchantID int64, paymentNo string) string {
	payload := fmt.Sprintf("%d|%d|%s", tenantID, merchantID, paymentNo)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
func ParseState(secret, raw string) (tenantID, merchantID int64, paymentNo string, err error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 2 {
		return 0, 0, "", errors.New("invalid paypal state")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, 0, "", err
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, 0, "", err
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return 0, 0, "", errors.New("invalid paypal state signature")
	}
	fields := strings.Split(string(payload), "|")
	if len(fields) != 3 {
		return 0, 0, "", errors.New("invalid paypal state")
	}
	tenantID, err = strconv.ParseInt(fields[0], 10, 64)
	if err != nil || tenantID <= 0 {
		return 0, 0, "", errors.New("invalid tenant in paypal state")
	}
	merchantID, err = strconv.ParseInt(fields[1], 10, 64)
	if err != nil || merchantID <= 0 {
		return 0, 0, "", errors.New("invalid merchant in paypal state")
	}
	paymentNo = fields[2]
	if paymentNo == "" {
		return 0, 0, "", errors.New("invalid payment in paypal state")
	}
	return tenantID, merchantID, paymentNo, nil
}
