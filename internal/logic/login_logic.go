package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.AuthTokenResp, error) {
	resp, err := l.svcCtx.UserClient.Login(l.ctx, &travel.LoginRequest{
		Username: req.Username, Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &types.AuthTokenResp{Token: resp.Token, User: toUserProfile(resp.User)}, nil
}
