package logic

import (
	"context"
	"errors"

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

	// If approved, trigger PayPal refund
	if req.Action == "APPROVE" && resp2.GetStatus() == "REFUNDED" {
		l.processPayPalRefund(orderNo)
	}

	return &types.HandleRefundResp{
		OrderNo:      resp2.GetOrderNo(),
		Status:       resp2.GetStatus(),
		RejectReason: resp2.GetRejectReason(),
	}, nil
}

func (l *HandleRefundLogic) processPayPalRefund(orderNo string) {
	// First find the order to get its ent ID
	o, err := l.svcCtx.DB.Order.Query().
		Where(order.OrderNoEQ(orderNo)).
		First(l.ctx)
	if err != nil {
		logx.Errorf("[PayPal Refund] failed to find order %s: %v", orderNo, err)
		return
	}
	// Then find the payment by order_id
	p, err := l.svcCtx.DB.Payment.Query().
		Where(payment.OrderIDEQ(int64(o.ID))).
		First(l.ctx)
	if err != nil {
		logx.Errorf("[PayPal Refund] failed to find payment for order %s: %v", orderNo, err)
		return
	}
	if p.ProviderPaymentID == "" || p.Provider != "paypal" {
		logx.Infof("[PayPal Refund] no provider payment id or not paypal, skipping")
		return
	}
	result, err := l.svcCtx.PayPal.RefundCapture(l.ctx, p.ProviderPaymentID, p.Amount, p.Currency)
	if err != nil {
		logx.Errorf("[PayPal Refund] refund failed for capture %s: %v", p.ProviderPaymentID, err)
		return
	}
	logx.Infof("[PayPal Refund] success: refund_id=%s status=%s", result.RefundID, result.Status)
}
