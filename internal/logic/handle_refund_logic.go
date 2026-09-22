package logic

import (
	"context"
	"errors"

	paymentregistry "gitee.com/meinongyihe/travel-api/internal/payment"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/order"
	"gitee.com/meinongyihe/travel-rpc/ent/payment"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleRefundLogic {
	return &HandleRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *HandleRefundLogic) HandleRefund(req *types.HandleRefundReq) (resp *types.HandleRefundResp, err error) {
	orderNo := pathVal(l.ctx, "orderNo")
	if orderNo == "" {
		return nil, errors.New("invalid order number")
	}

	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}

	resp2, err := l.svcCtx.OrderClient.HandleRefund(ctx, &travel.HandleRefundRequest{
		OrderNo: orderNo,
		Action:  req.Action,
		Reason:  req.Reason,
	})
	if err != nil {
		return nil, err
	}

	// If approved, trigger provider-specific refund based on payment provider.
	if req.Action == "APPROVE" && resp2.GetStatus() == "REFUNDED" {
		l.processRefund(orderNo)
	}

	return &types.HandleRefundResp{
		OrderNo:      resp2.GetOrderNo(),
		Status:       resp2.GetStatus(),
		RejectReason: resp2.GetRejectReason(),
	}, nil
}

// processRefund routes the refund to the correct payment provider based on the provider field.
func (l *HandleRefundLogic) processRefund(orderNo string) {
	// Find the order to get its ent ID.
	o, err := l.svcCtx.DB.Order.Query().
		Where(order.OrderNoEQ(orderNo)).
		First(l.ctx)
	if err != nil {
		logx.Errorf("[Refund] failed to find order %s: %v", orderNo, err)
		return
	}
	// Find the payment by order_id.
	p, err := l.svcCtx.DB.Payment.Query().
		Where(payment.OrderIDEQ(int64(o.ID))).
		First(l.ctx)
	if err != nil {
		logx.Errorf("[Refund] failed to find payment for order %s: %v", orderNo, err)
		return
	}
	if p.ProviderPaymentID == "" {
		logx.Infof("[Refund] no provider payment id for order %s, skipping", orderNo)
		return
	}

	// Look up the provider from the registry.
	provider, ok := paymentregistry.Get(p.Provider)
	if !ok {
		logx.Errorf("[Refund] unknown payment provider %q for order %s", p.Provider, orderNo)
		return
	}

	result, err := provider.Refund(l.ctx, p.ProviderPaymentID, p.Amount, p.Currency)
	if err != nil {
		logx.Errorf("[Refund] %s refund failed for payment %s: %v", p.Provider, p.ProviderPaymentID, err)
		return
	}
	logx.Infof("[Refund] %s success: refund_id=%s status=%s", p.Provider, result.RefundID, result.Status)
}
