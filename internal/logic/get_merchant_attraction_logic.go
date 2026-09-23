package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/attraction"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMerchantAttractionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantAttractionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantAttractionLogic {
	return &GetMerchantAttractionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetMerchantAttractionLogic) GetMerchantAttraction() (resp *types.Attraction, err error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	merchantID, _ := middleware.MerchantID(l.ctx)

	a, err := l.svcCtx.DB.Attraction.Query().
		Where(attraction.IDEQ(int(id)), attraction.MerchantIDEQ(merchantID)).
		Only(l.ctx)
	if err != nil {
		return nil, errors.New("attraction not found")
	}
	out := toAttraction(a)
	return &out, nil
}
