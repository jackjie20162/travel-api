package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-api/internal/utils"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/pkg/errors"
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
	// 1. 校验图形验证码
	if !utils.VerifyCaptcha(l.svcCtx.Redis, l.svcCtx.Config.Captcha.KeyPrefix, req.CaptchaId, req.CaptchaAnswer) {
		return nil, errors.New("图形验证码错误或已过期")
	}

	// 2. 校验邮件验证码
	if !utils.VerifyEmailCode(l.svcCtx.Redis, req.Email, req.EmailCode) {
		return nil, errors.New("邮件验证码错误或已过期")
	}

	// 3. 调用 RPC 邮箱登录（无密码）
	resp, err := l.svcCtx.UserClient.Login(l.ctx, &travel.LoginRequest{
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &types.AuthTokenResp{Token: resp.Token, User: toUserProfile(resp.User)}, nil
}
