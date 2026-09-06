package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
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

	return &types.HandleRefundResp{
		OrderNo:      resp2.GetOrderNo(),
		Status:       resp2.GetStatus(),
		RejectReason: resp2.GetRejectReason(),
	}, nil
}
