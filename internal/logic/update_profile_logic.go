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

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileReq) (*types.UserProfile, error) {
	userID, ok := middleware.AuthUserID(l.ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	ctx, err := userRpcContext(l.ctx, userID)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.UserClient.UpdateProfile(ctx, &travel.UpdateProfileRequest{
		Nickname: req.Nickname, Avatar: req.Avatar,
		Email: req.Email, Mobile: req.Mobile, Bio: req.Bio,
		Birthday: req.Birthday, Gender: req.Gender,
	})
	if err != nil {
		return nil, err
	}
	out := toUserProfile(resp)
	return &out, nil
}
