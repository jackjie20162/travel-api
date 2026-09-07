package logic

import (
	"context"
	"fmt"

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
	// 1. 校验图形验证码
	if !utils.VerifyCaptcha(l.svcCtx.Redis, l.svcCtx.Config.Captcha.KeyPrefix, req.CaptchaId, req.CaptchaAnswer) {
		return nil, errors.New("图形验证码错误或已过期")
	}

	// 2. 生成邮件验证码
	code, err := utils.GenerateEmailCode()
	if err != nil {
		return nil, errors.Wrap(err, "生成验证码失败")
	}

	// 3. 存入 Redis
	if err := utils.StoreEmailCode(l.svcCtx.Redis, req.Email, code, l.svcCtx.Config.Captcha.EmailCode); err != nil {
		return nil, errors.Wrap(err, "存储验证码失败")
	}

	// 4. 发送邮件
	subject := fmt.Sprintf("验证码：%s", code)
	body := fmt.Sprintf(`<div style="font-family:sans-serif;padding:20px">
		<h2>您的验证码</h2>
		<p>验证码: <strong style="font-size:24px;color:#e74c3c">%s</strong></p>
		<p>有效期 %d 秒，请勿泄露给他人。</p>
	</div>`, code, l.svcCtx.Config.Captcha.EmailCode)

	if err := utils.Email(l.svcCtx.Config.Email, req.Email, subject, body); err != nil {
		return nil, errors.Wrap(err, "发送邮件失败")
	}

	return &types.OkResp{Ok: true, Msg: "验证码已发送"}, nil
}
