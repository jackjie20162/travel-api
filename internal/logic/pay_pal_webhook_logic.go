package logic

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayPalWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayPalWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayPalWebhookLogic {
	return &PayPalWebhookLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PayPalWebhookLogic) PayPalWebhook(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	// Verify webhook signature
	err = l.svcCtx.PayPal.VerifyWebhook(l.ctx, r.Header, json.RawMessage(body))
	if err != nil {
		return err
	}
	// Parse event
	var event struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	// Handle CHECKOUT.ORDER.APPROVED and PAYMENT.CAPTURE.COMPLETED
	if event.EventType == "PAYMENT.CAPTURE.COMPLETED" || event.EventType == "CHECKOUT.ORDER.APPROVED" {
		// Find payment by provider payment ID and mark as paid
		ctx := publicRpcContext(l.ctx)
		_, _ = l.svcCtx.PaymentClient.MarkPaid(ctx, &travel.MarkPaymentPaidRequest{
			ProviderPaymentId: event.Resource.ID,
		})
	}
	return nil
}
