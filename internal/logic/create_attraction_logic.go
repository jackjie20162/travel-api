package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAttractionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAttractionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAttractionLogic {
	return &CreateAttractionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *CreateAttractionLogic) CreateAttraction(req *types.AttractionReq) (resp *types.Attraction, err error) {
	tenantID, _ := middleware.TenantID(l.ctx)
	merchantID, _ := middleware.MerchantID(l.ctx)
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	a, err := l.svcCtx.DB.Attraction.Create().
		SetTenantID(tenantID).
		SetMerchantID(merchantID).
		SetName(req.Name).
		SetNameEn(req.NameEn).
		SetCategory(req.Category).
		SetTags(req.Tags).
		SetCity(req.City).
		SetCityID(req.CityId).
		SetAddress(req.Address).
		SetLatitude(req.Latitude).
		SetLongitude(req.Longitude).
		SetPoiID(req.PoiId).
		SetPoiName(req.PoiName).
		SetCoverImage(req.CoverImage).
		SetImages(req.Images).
		SetDescription(req.Description).
		SetRichContent(req.RichContent).
		SetOpeningHours(req.OpeningHours).
		SetSuggestedDuration(req.SuggestedDuration).
		SetTicketInfo(req.TicketInfo).
		SetSort(req.Sort).
		SetStatus(status).
		Save(l.ctx)
	if err != nil {
		return nil, err
	}
	// 写入多语言译文（en-US/ar），基础语言内容已在主表字段。
	if err := saveAttractionTranslations(l.ctx, l.svcCtx.DB, int64(a.ID), req.Translations); err != nil {
		return nil, err
	}
	out := toAttraction(a)
	items := []types.Attraction{out}
	if err := attachAttractionTranslations(l.ctx, l.svcCtx.DB, items); err != nil {
		return nil, err
	}
	out = items[0]
	return &out, nil
}
