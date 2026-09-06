package logic

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/paypal"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type CreatePaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreatePaymentLogic) CreatePayment(req *types.CreatePaymentReq) (*types.PaymentResp, error) {
	if req.OrderNo == "" || req.Provider == "" || req.IdempotencyKey == "" {
		return nil, errors.New("orderNo, provider and idempotencyKey are required")
	}
	if !strings.EqualFold(req.Provider, "paypal") {
		return nil, errors.New("only paypal is enabled")
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	p, err := l.svcCtx.PaymentClient.Create(ctx, &travel.CreatePaymentRequest{
		OrderNo: req.OrderNo, Provider: req.Provider, IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	tenantID, ok := middleware.TenantID(l.ctx)
	if !ok {
		return nil, errors.New("tenant context is required")
	}
	merchantID, ok := middleware.MerchantID(l.ctx)
	if !ok {
		return nil, errors.New("merchant context is required")
	}
	secret := l.paypalSecret()
	if secret == "" {
		return nil, errors.New("paypal client secret is not configured")
	}
	state := paypal.SignState(secret, tenantID, merchantID, p.PaymentNo)
	checkout, err := l.svcCtx.PayPal.CreateOrder(ctx, paypal.CheckoutInput{
		PaymentNo: p.PaymentNo, Amount: p.Amount, Currency: p.Currency, State: state,
	})
	if err != nil {
		return nil, err
	}
	p, err = l.svcCtx.PaymentClient.SetProviderID(ctx, &travel.SetPaymentProviderIDRequest{
		PaymentNo: p.PaymentNo, ProviderPaymentId: checkout.OrderID,
	})
	if err != nil {
		return nil, err
	}
	out := toPaymentResp(p)
	out.CheckoutURL = checkout.CheckoutURL
	return out, nil
}

func (l *CreatePaymentLogic) paypalSecret() string {
	if l.svcCtx.Config.Payment.PayPal.Secret != "" {
		return l.svcCtx.Config.Payment.PayPal.Secret
	}
	return os.Getenv("PAYPAL_CLIENT_SECRET")
}

func paymentContext(ctx context.Context, tenantID, merchantID int64) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"x-tenant-id", strconv.FormatInt(tenantID, 10),
		"x-merchant-id", strconv.FormatInt(merchantID, 10),
	))
}
