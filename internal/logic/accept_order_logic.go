package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type AcceptOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAcceptOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptOrderLogic {
	return &AcceptOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *AcceptOrderLogic) AcceptOrder(req *types.AcceptOrderReq) (resp *types.AcceptOrderResp, err error) {
	orderNo := pathVal(l.ctx, "orderNo")
	if orderNo == "" {
		return nil, errors.New("invalid order number")
	}

	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}

	resp2, err := l.svcCtx.OrderClient.AcceptOrder(ctx, &travel.AcceptOrderRequest{
		OrderNo: orderNo,
		Action:  req.Action,
		Reason:  req.Reason,
	})
	if err != nil {
		return nil, err
	}

	return &types.AcceptOrderResp{
		OrderNo:      resp2.GetOrderNo(),
		Status:       resp2.GetStatus(),
		RejectReason: resp2.GetRejectReason(),
	}, nil
}
