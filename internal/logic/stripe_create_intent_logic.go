package logic

import (
	"context"
	"errors"
	"strconv"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/payment"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type StripeCreateIntentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStripeCreateIntentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StripeCreateIntentLogic {
	return &StripeCreateIntentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// StripeCreateIntent creates a Stripe PaymentIntent for embedded Elements checkout.
// Flow: validate → RPC create payment (provider=stripe) → Stripe API create PaymentIntent → return client_secret.
func (l *StripeCreateIntentLogic) StripeCreateIntent(req *types.CreatePaymentReq) (*types.StripeIntentResp, error) {
	if req.OrderNo == "" || req.IdempotencyKey == "" {
		return nil, errors.New("orderNo and idempotencyKey are required")
	}
	if l.svcCtx.StripeProvider == nil {
		return nil, errors.New("stripe payment is not configured")
	}

	// 1. Create payment record via RPC (provider="stripe").
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	p, err := l.svcCtx.PaymentClient.Create(ctx, &travel.CreatePaymentRequest{
		OrderNo:        req.OrderNo,
		Provider:       "stripe",
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}

	// 2. Get tenant/merchant for metadata.
	tenantID, _ := middleware.TenantID(l.ctx)
	merchantID, _ := middleware.MerchantID(l.ctx)

	// 3. Call Stripe provider to create PaymentIntent.
	result, err := l.svcCtx.StripeProvider.CreateCheckout(l.ctx, checkoutInputFromPayment(p, req.OrderNo, tenantID, merchantID))
	if err != nil {
		return nil, err
	}

	// 4. Store the Stripe PaymentIntent ID on the payment record.
	p, err = l.svcCtx.PaymentClient.SetProviderID(ctx, &travel.SetPaymentProviderIDRequest{
		PaymentNo:         p.PaymentNo,
		ProviderPaymentId: result.ProviderPaymentID,
	})
	if err != nil {
		return nil, err
	}

	return &types.StripeIntentResp{
		PaymentNo:      p.PaymentNo,
		ClientSecret:   result.ClientSecret,
		PublishableKey: l.svcCtx.StripeProvider.PublishableKey(),
	}, nil
}

func checkoutInputFromPayment(p *travel.Payment, orderNo string, tenantID, merchantID int64) payment.CheckoutInput {
	return payment.CheckoutInput{
		PaymentNo: p.PaymentNo,
		OrderNo:   orderNo,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Metadata: map[string]string{
			"tenant_id":   strconv.FormatInt(tenantID, 10),
			"merchant_id": strconv.FormatInt(merchantID, 10),
		},
	}
}
