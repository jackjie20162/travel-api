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

type ShopifyWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShopifyWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShopifyWebhookLogic {
	return &ShopifyWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// ShopifyWebhook handles incoming Shopify webhook events (orders/create, orders/paid, orders/cancelled).
func (l *ShopifyWebhookLogic) ShopifyWebhook(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	provider, ok := paymentregistry.Get("shopify")
	if !ok {
		return errors.New("shopify provider not registered")
	}

	evt, err := provider.VerifyWebhook(l.ctx, r.Header, body)
	if err != nil {
		return err
	}

	rpcCtx := publicRpcContext(l.ctx)

	switch evt.Type {
	case paymentregistry.EventPaymentSucceeded:
		// orders/paid — mark the corresponding travel payment as PAID.
		if evt.PaymentNo != "" {
			_, err = l.svcCtx.PaymentClient.MarkPaid(rpcCtx, &travel.MarkPaymentPaidRequest{
				PaymentNo:         evt.PaymentNo,
				ProviderPaymentId: evt.ProviderPaymentID,
			})
			if err != nil {
				logx.Errorf("[Shopify Webhook] MarkPaid failed for payment_no=%s: %v", evt.PaymentNo, err)
				return err
			}
			l.Infof("[Shopify Webhook] payment marked PAID: payment_no=%s shopify_order=%s", evt.PaymentNo, evt.ProviderPaymentID)
		}

	case paymentregistry.EventOrderCreated:
		// orders/create — log for now; full order mapping can be implemented later.
		l.Infof("[Shopify Webhook] order created: shopify_order=%s payment_no=%s", evt.ProviderPaymentID, evt.PaymentNo)

	case paymentregistry.EventOrderCancelled:
		l.Infof("[Shopify Webhook] order cancelled: shopify_order=%s", evt.ProviderPaymentID)
	}

	return nil
}
