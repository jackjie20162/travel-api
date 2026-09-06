package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *CancelOrderLogic) CancelOrder() (resp *types.CancelOrderResp, err error) {
	orderNo := pathVal(l.ctx, "orderNo")
	if orderNo == "" {
		return nil, errors.New("invalid order number")
	}

	userID, ok := middleware.AuthUserID(l.ctx)
	if !ok || userID <= 0 {
		return nil, errors.New("missing user context")
	}

	ctx, err := userRpcContext(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	resp2, err := l.svcCtx.OrderClient.CancelOrder(ctx, &travel.OrderNoRequest{
		OrderNo: orderNo,
	})
	if err != nil {
		return nil, err
	}

	return &types.CancelOrderResp{
		OrderNo: resp2.GetOrderNo(),
		Status:  resp2.GetStatus(),
	}, nil
}
