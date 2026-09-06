package logic

import (
	"context"
	"errors"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentLogic {
	return &GetPaymentLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetPaymentLogic) GetPayment() (*types.PaymentResp, error) {
	paymentNo := pathVal(l.ctx, "paymentNo")
	if paymentNo == "" {
		return nil, errors.New("invalid payment number")
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	p, err := l.svcCtx.PaymentClient.Get(ctx, &travel.PaymentNoRequest{PaymentNo: paymentNo})
	if err != nil {
		return nil, err
	}
	return toPaymentResp(p), nil
}
