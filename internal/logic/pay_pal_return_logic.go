package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"

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

func (l *PayPalReturnLogic) PayPalReturn() (*types.PaymentResp, error) {
	// The handler extracts query params and calls this logic
	// For now, return error as the handler handles PayPal return directly
	return nil, errors.New("paypal return handled in handler")
}
