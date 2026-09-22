package logic

import (
	"context"
	"errors"
	"io"
	"net/http"

	paymentregistry "gitee.com/meinongyihe/travel-api/internal/payment"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type StripeWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStripeWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StripeWebhookLogic {
	return &StripeWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// StripeWebhook handles incoming Stripe webhook events.
// It verifies the signature, extracts the normalised event, and calls MarkPaid on success.
func (l *StripeWebhookLogic) StripeWebhook(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Look up the Stripe provider from the registry.
	provider, ok := paymentregistry.Get("stripe")
	if !ok {
		return errors.New("stripe provider not registered")
	}

	// Verify signature and parse the event.
	evt, err := provider.VerifyWebhook(l.ctx, r.Header, body)
	if err != nil {
		return err
	}

	// Handle payment success.
	if evt.Type == paymentregistry.EventPaymentSucceeded && evt.PaymentNo != "" {
		rpcCtx := publicRpcContext(l.ctx)
		_, err = l.svcCtx.PaymentClient.MarkPaid(rpcCtx, &travel.MarkPaymentPaidRequest{
			PaymentNo:         evt.PaymentNo,
			ProviderPaymentId: evt.ProviderPaymentID,
		})
		if err != nil {
			// Idempotent: if already PAID, MarkPaid returns success.
			logx.Errorf("[Stripe Webhook] MarkPaid failed for payment_no=%s: %v", evt.PaymentNo, err)
			return err
		}
		l.Infof("[Stripe Webhook] payment marked PAID: payment_no=%s pi=%s", evt.PaymentNo, evt.ProviderPaymentID)
	}

	return nil
}
