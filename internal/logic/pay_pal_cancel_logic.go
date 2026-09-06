package logic

import (
	"context"
	"errors"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type PayPalCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayPalCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayPalCancelLogic {
	return &PayPalCancelLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PayPalCancelLogic) PayPalCancel() (*types.OkResp, error) {
	return nil, errors.New("paypal cancel handled in handler")
}
