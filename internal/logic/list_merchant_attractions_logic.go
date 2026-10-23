package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/attraction"
	"gitee.com/meinongyihe/travel-rpc/ent/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMerchantAttractionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMerchantAttractionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMerchantAttractionsLogic {
	return &ListMerchantAttractionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// toAttraction 将 ent.Attraction 映射为 API 响应类型。
func toAttraction(a *ent.Attraction) types.Attraction {
	return types.Attraction{
		Id: int64(a.ID), TenantId: a.TenantID, MerchantId: a.MerchantID,
		Name: a.Name, NameEn: a.NameEn, Category: a.Category, Tags: a.Tags,
		City: a.City, CityId: a.CityID, Address: a.Address, Latitude: a.Latitude, Longitude: a.Longitude,
		PoiId: a.PoiID, PoiName: a.PoiName, CoverImage: a.CoverImage, Images: a.Images,
		Description: a.Description, RichContent: a.RichContent, OpeningHours: a.OpeningHours,
		SuggestedDuration: a.SuggestedDuration, TicketInfo: a.TicketInfo,
		Sort: a.Sort, Status: a.Status,
	}
}

// queryAttractions 通用景点查询：按商户隔离 + 过滤 + 分页（商户端与公开端共用）。
func queryAttractions(ctx context.Context, db *ent.Client, req *types.AttractionListReq, merchantID int64) (*types.AttractionListResp, error) {
	preds := []predicate.Attraction{attraction.MerchantIDEQ(merchantID)}
	if req.Keyword != "" {
		preds = append(preds, attraction.Or(
			attraction.NameContains(req.Keyword),
			attraction.NameEnContains(req.Keyword),
			attraction.CityContains(req.Keyword),
		))
	}
	if req.Category != "" {
		preds = append(preds, attraction.CategoryEQ(req.Category))
	}
	if req.Status != "" {
		preds = append(preds, attraction.StatusEQ(req.Status))
	}

	q := db.Attraction.Query().Where(preds...)
	total, err := q.Count(ctx)
	if err != nil {
		return nil, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	rows, err := q.Order(attraction.BySort(), attraction.ByID()).
		Offset((page - 1) * pageSize).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]types.Attraction, 0, len(rows))
	for _, a := range rows {
		items = append(items, toAttraction(a))
	}
	return &types.AttractionListResp{Items: items, Total: int64(total)}, nil
}

func (l *ListMerchantAttractionsLogic) ListMerchantAttractions(req *types.AttractionListReq) (resp *types.AttractionListResp, err error) {
	merchantID, _ := middleware.MerchantID(l.ctx)
	resp, err = queryAttractions(l.ctx, l.svcCtx.DB, req, merchantID)
	if err != nil {
		return nil, err
	}
	// 管理端需要全部语种译文明细用于编辑。
	if err := attachAttractionTranslations(l.ctx, l.svcCtx.DB, resp.Items); err != nil {
		return nil, err
	}
	return resp, nil
}
