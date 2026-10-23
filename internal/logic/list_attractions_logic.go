package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/i18n"
	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAttractionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAttractionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAttractionsLogic {
	return &ListAttractionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// ListAttractions 公开景点列表：默认仅返回 ACTIVE 状态景点。
func (l *ListAttractionsLogic) ListAttractions(req *types.AttractionListReq) (resp *types.AttractionListResp, err error) {
	merchantID, _ := middleware.MerchantID(l.ctx)
	if req.Status == "" {
		req.Status = "ACTIVE"
	}
	resp, err = queryAttractions(l.ctx, l.svcCtx.DB, req, merchantID)
	if err != nil {
		return nil, err
	}
	// 公开端按请求语言解析展示，命中译文则覆盖，否则回退基础语言。
	locale := req.Locale
	if locale == "" {
		locale = i18n.LocaleFrom(l.ctx)
	}
	if err := resolveAttractionLocale(l.ctx, l.svcCtx.DB, resp.Items, locale); err != nil {
		return nil, err
	}
	return resp, nil
}
