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

type UpdateAttractionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAttractionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAttractionLogic {
	return &UpdateAttractionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *UpdateAttractionLogic) UpdateAttraction(req *types.AttractionReq) (resp *types.Attraction, err error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	merchantID, _ := middleware.MerchantID(l.ctx)

	exist, err := l.svcCtx.DB.Attraction.Query().
		Where(attraction.IDEQ(int(id)), attraction.MerchantIDEQ(merchantID)).
		Exist(l.ctx)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("attraction not found")
	}

	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	a, err := l.svcCtx.DB.Attraction.UpdateOneID(int(id)).
		SetName(req.Name).
		SetNameEn(req.NameEn).
		SetCategory(req.Category).
		SetTags(req.Tags).
		SetCity(req.City).
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
	out := toAttraction(a)
	return &out, nil
}
