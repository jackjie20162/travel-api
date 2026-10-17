package logic

import (
	"context"
	"encoding/json"
	"strconv"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/itinerarystop"
	"gitee.com/meinongyihe/travel-rpc/ent/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAttractionProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAttractionProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAttractionProductsLogic {
	return &ListAttractionProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// ListAttractionProducts 公开接口：返回行程中包含指定景点的已发布线路（产品）。
// 绑定关系记录在 itinerary_stop.type_params JSON 的 attractionId 字段中，
// 先用 LIKE 粗筛再精确解析，按产品去重、最低价升序返回。
func (l *ListAttractionProductsLogic) ListAttractionProducts() (resp *types.ProductListResp, err error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	merchantID, _ := middleware.MerchantID(l.ctx)

	stops, err := l.svcCtx.DB.ItineraryStop.Query().
		Where(
			itinerarystop.MerchantIDEQ(merchantID),
			itinerarystop.Or(
				itinerarystop.TypeParamsContains(`"attractionId":`+strconv.FormatInt(id, 10)),
				itinerarystop.TypeParamsContains(`"attractionId":"`+strconv.FormatInt(id, 10)+`"`),
			),
		).
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	seen := make(map[int64]bool)
	ids := make([]int, 0)
	for _, s := range stops {
		if s.ProductID == 0 || seen[s.ProductID] {
			continue
		}
		var tp map[string]any
		if json.Unmarshal([]byte(s.TypeParams), &tp) != nil {
			continue
		}
		if parseAttractionID(tp["attractionId"]) != int(id) {
			continue
		}
		seen[s.ProductID] = true
		ids = append(ids, int(s.ProductID))
	}

	empty := &types.ProductListResp{Items: []types.Product{}, Total: 0}
	if len(ids) == 0 {
		return empty, nil
	}

	rows, err := l.svcCtx.DB.Product.Query().
		Where(product.IDIn(ids...), product.StatusEQ("PUBLISHED")).
		Order(product.ByMinPrice(), product.ByID()).
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	items := make([]types.Product, 0, len(rows))
	for _, p := range rows {
		items = append(items, entToProduct(p))
	}
	// 线路没有主图时回退取行程绑定景点的图
	fillMissingProductCovers(l.ctx, l.svcCtx.DB, items)
	return &types.ProductListResp{Items: items, Total: int64(len(items))}, nil
}

// entToProduct 将 ent.Product 映射为 API 响应类型（公开列表口径）。
func entToProduct(p *ent.Product) types.Product {
	return types.Product{
		Id: int64(p.ID), TenantId: p.TenantID, MerchantId: p.MerchantID,
		Code: p.Code, Title: p.Title, Slug: p.Slug,
		Destination: p.Destination, Description: p.Description,
		Currency: p.Currency, MinPrice: p.MinPrice, Status: p.Status,
		Highlights: p.Highlights, CoverImage: p.CoverImage,
		Images: p.Images, VideoUrl: p.VideoURL,
		RichContent: p.RichContent, BookingNotice: p.BookingNotice,
	}
}
