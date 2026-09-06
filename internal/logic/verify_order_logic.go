package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyOrderLogic {
	return &VerifyOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *VerifyOrderLogic) VerifyOrder(req *types.VerifyOrderReq) (*types.VerifyOrderResp, error) {
	orderNo := pathVal(l.ctx, "orderNo")
	if orderNo == "" {
		return nil, errors.New("invalid order number")
	}

	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.OrderClient.VerifyOrder(ctx, &travel.VerifyOrderRequest{
		OrderNo: orderNo,
		Action:  req.Action,
		Reason:  req.Reason,
	})
	if err != nil {
		return nil, err
	}

	return &types.VerifyOrderResp{
		OrderNo:      resp.GetOrderNo(),
		Status:       resp.GetStatus(),
		RejectReason: resp.GetRejectReason(),
		VerifiedAt:   resp.GetVerifiedAt(),
	}, nil
}
