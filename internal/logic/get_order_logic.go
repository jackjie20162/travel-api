package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitee.com/meinongyihe/travel-rpc/travel"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetOrderLogic) GetOrder() (*types.CreateOrderResp, error) {
	orderNo := pathVal(l.ctx, "orderNo")
	if orderNo == "" {
		return nil, errors.New("invalid order number")
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.OrderClient.Get(ctx, &travel.OrderNoRequest{OrderNo: orderNo})
	if err != nil {
		return nil, err
	}
	return toOrderResp(resp), nil
}
