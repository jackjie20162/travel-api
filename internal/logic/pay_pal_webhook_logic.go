package logic

import (
	"context"
	"errors"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PayPalWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayPalWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayPalWebhookLogic {
	return &PayPalWebhookLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PayPalWebhookLogic) PayPalWebhook() error {
	return errors.New("paypal webhook handled in handler")
}
