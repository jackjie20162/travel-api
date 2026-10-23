package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/i18n"
	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/attraction"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAttractionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAttractionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAttractionLogic {
	return &GetAttractionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// GetAttraction 公开景点详情：仅返回 ACTIVE 状态景点，与公开列表口径一致。
func (l *GetAttractionLogic) GetAttraction() (resp *types.Attraction, err error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	merchantID, _ := middleware.MerchantID(l.ctx)

	a, err := l.svcCtx.DB.Attraction.Query().
		Where(attraction.IDEQ(int(id)), attraction.MerchantIDEQ(merchantID), attraction.StatusEQ("ACTIVE")).
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("attraction not found")
		}
		return nil, err
	}
	out := toAttraction(a)
	items := []types.Attraction{out}
	if err := resolveAttractionLocale(l.ctx, l.svcCtx.DB, items, i18n.LocaleFrom(l.ctx)); err != nil {
		return nil, err
	}
	out = items[0]
	return &out, nil
}
