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

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProfileLogic) GetProfile() (*types.UserProfile, error) {
	userID, ok := middleware.AuthUserID(l.ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	resp, err := l.svcCtx.UserClient.GetProfile(l.ctx, &travel.UserIdRequest{Id: userID})
	if err != nil {
		return nil, err
	}
	out := toUserProfile(resp)
	return &out, nil
}
