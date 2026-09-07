package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-api/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaLogic {
	return &CaptchaLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CaptchaLogic) Captcha() (*types.CaptchaResp, error) {
	id, svg, err := utils.GenerateCaptcha(l.svcCtx.Redis, l.svcCtx.Config.Captcha.KeyPrefix, l.svcCtx.Config.Captcha.Expire)
	if err != nil {
		return nil, err
	}
	return &types.CaptchaResp{
		CaptchaId:    id,
		CaptchaImage: svg,
	}, nil
}
