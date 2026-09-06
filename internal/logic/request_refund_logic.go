package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"gitee.com/meinongyihe/travel-api/internal/middleware"

	"github.com/zeromicro/go-zero/core/logx"
)

type RequestRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRequestRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RequestRefundLogic {
	return &RequestRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *RequestRefundLogic) RequestRefund(req *types.RequestRefundReq) (resp *types.RequestRefundResp, err error) {
	orderNo := pathVal(l.ctx, "orderNo")
	if orderNo == "" {
		return nil, errors.New("invalid order number")
	}

	// Get user ID from middleware for customer operations
	userID, ok := middleware.AuthUserID(l.ctx)
	if !ok || userID <= 0 {
		return nil, errors.New("missing user context")
	}

	ctx, err := userRpcContext(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	resp2, err := l.svcCtx.OrderClient.RequestRefund(ctx, &travel.RefundRequest{
		OrderNo: orderNo,
		Reason:  req.Reason,
	})
	if err != nil {
		return nil, err
	}

	return &types.RequestRefundResp{
		OrderNo:      resp2.GetOrderNo(),
		Status:       resp2.GetStatus(),
		RejectReason: resp2.GetRejectReason(),
	}, nil
}
