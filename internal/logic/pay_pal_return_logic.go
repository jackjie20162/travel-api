package logic

import (
	"context"
	"errors"
	"net/http"

	"gitee.com/meinongyihe/travel-api/internal/paypal"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayPalReturnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayPalReturnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayPalReturnLogic {
	return &PayPalReturnLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PayPalReturnLogic) PayPalReturn(r *http.Request) (*types.PaymentResp, error) {
	state := r.URL.Query().Get("state")
	if state == "" {
		return nil, errors.New("missing state parameter")
	}
	secret := l.svcCtx.Config.Payment.PayPal.Secret
	tenantID, merchantID, paymentNo, err := paypal.ParseState(secret, state)
	if err != nil {
		return nil, err
	}
	ctx := paymentContext(l.ctx, tenantID, merchantID)
	// Get payment to check current status
	p, err := l.svcCtx.PaymentClient.Get(ctx, &travel.PaymentNoRequest{PaymentNo: paymentNo})
	if err != nil {
		return nil, err
	}
	// If already paid (e.g. webhook already processed), return directly
	if p.Status == "PAID" {
		return toPaymentResp(p), nil
	}
	// Capture PayPal order
	capture, err := l.svcCtx.PayPal.CaptureOrder(l.ctx, p.ProviderPaymentId)
	if err != nil {
		// If capture fails because it's already captured, check payment status again
		p2, getErr := l.svcCtx.PaymentClient.Get(ctx, &travel.PaymentNoRequest{PaymentNo: paymentNo})
		if getErr == nil && p2.Status == "PAID" {
			return toPaymentResp(p2), nil
		}
		return nil, err
	}
	// Mark payment as paid
	paid, err := l.svcCtx.PaymentClient.MarkPaid(ctx, &travel.MarkPaymentPaidRequest{
		PaymentNo:         paymentNo,
		ProviderPaymentId: capture.CaptureID,
	})
	if err != nil {
		// If MarkPaid fails, check if it was already marked paid by webhook
		p2, getErr := l.svcCtx.PaymentClient.Get(ctx, &travel.PaymentNoRequest{PaymentNo: paymentNo})
		if getErr == nil && p2.Status == "PAID" {
			return toPaymentResp(p2), nil
		}
		return nil, err
	}
	out := toPaymentResp(paid)
	return out, nil
}
