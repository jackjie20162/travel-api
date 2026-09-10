package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/i18n"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-api/internal/utils"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendEmailCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailCodeLogic {
	return &SendEmailCodeLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SendEmailCodeLogic) SendEmailCode(req *types.SendEmailCodeReq) (*types.OkResp, error) {
	locale := i18n.LocaleFrom(l.ctx)

	// 1. 校验图形验证码
	if !utils.VerifyCaptcha(l.svcCtx.Redis, l.svcCtx.Config.Captcha.KeyPrefix, req.CaptchaId, req.CaptchaAnswer) {
		return nil, errors.New("CAPTCHA_INVALID")
	}

	// 2. 生成邮件验证码
	code, err := utils.GenerateEmailCode()
	if err != nil {
		l.Errorf("generate email code failed: %v", err)
		return nil, errors.New("GENERATE_CODE_FAILED")
	}

	// 3. 存入 Redis
	if err := utils.StoreEmailCode(l.svcCtx.Redis, req.Email, code, l.svcCtx.Config.Captcha.EmailCode); err != nil {
		l.Errorf("store email code failed: %v", err)
		return nil, errors.New("STORE_CODE_FAILED")
	}

	// 4. 发送邮件（按 locale 选择三语模板）
	subject, body := i18n.EmailVerificationCode(locale, code, l.svcCtx.Config.Captcha.EmailCode)
	if err := utils.Email(l.svcCtx.Config.Email, req.Email, subject, body); err != nil {
		l.Errorf("send email failed: %v", err)
		return nil, errors.New("SEND_EMAIL_FAILED")
	}

	return &types.OkResp{Ok: true, Msg: i18n.T(locale, "EMAIL_CODE_SENT")}, nil
}
