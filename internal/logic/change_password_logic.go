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

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (*types.OkResp, error) {
	userID, ok := middleware.AuthUserID(l.ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	ctx, err := userRpcContext(l.ctx, userID)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.UserClient.ChangePassword(ctx, &travel.ChangePasswordRequest{
		OldPassword: req.OldPassword, NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	return &types.OkResp{Ok: true, Msg: "password changed"}, nil
}
