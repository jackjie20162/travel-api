package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByMobileLogic {
	return &LoginByMobileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *LoginByMobileLogic) LoginByMobile(req *types.LoginByMobileReq) (*types.AuthTokenResp, error) {
	resp, err := l.svcCtx.UserClient.LoginByMobile(l.ctx, &travel.LoginByMobileRequest{
		Mobile: req.Mobile, Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &types.AuthTokenResp{Token: resp.Token, User: toUserProfile(resp.User)}, nil
}
