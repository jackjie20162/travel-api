package logic

import (
	"context"
	"encoding/json"
	"strconv"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/attraction"
	"gitee.com/meinongyihe/travel-rpc/ent/itinerarystop"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductAttractionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductAttractionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductAttractionsLogic {
	return &ListProductAttractionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// ListProductAttractions 返回某旅游路线（产品）行程节点绑定的景点。
// 绑定关系记录在 itinerary_stop.type_params JSON 的 attractionId 字段中，
// 按行程顺序去重后查询景点库返回。
func (l *ListProductAttractionsLogic) ListProductAttractions() (resp *types.AttractionListResp, err error) {
	productID, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}

	stops, err := l.svcCtx.DB.ItineraryStop.Query().
		Where(itinerarystop.ProductIDEQ(productID)).
		Order(itinerarystop.BySequence(), itinerarystop.ByID()).
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	seen := make(map[int]bool)
	order := make([]int, 0)
	for _, s := range stops {
		if s.TypeParams == "" {
			continue
		}
		var tp map[string]any
		if json.Unmarshal([]byte(s.TypeParams), &tp) != nil {
			continue
		}
		aid := parseAttractionID(tp["attractionId"])
		if aid > 0 && !seen[aid] {
			seen[aid] = true
			order = append(order, aid)
		}
	}

	empty := &types.AttractionListResp{Items: []types.Attraction{}, Total: 0}
	if len(order) == 0 {
		return empty, nil
	}

	rows, err := l.svcCtx.DB.Attraction.Query().
		Where(attraction.IDIn(order...)).
		All(l.ctx)
	if err != nil {
		return nil, err
	}
	byID := make(map[int]*ent.Attraction, len(rows))
	for _, a := range rows {
		byID[a.ID] = a
	}

	items := make([]types.Attraction, 0, len(order))
	for _, id := range order {
		if a, ok := byID[id]; ok {
			items = append(items, toAttraction(a))
		}
	}
	return &types.AttractionListResp{Items: items, Total: int64(len(items))}, nil
}

// parseAttractionID 兼容 JSON 数字(float64)、字符串两种存储形式。
func parseAttractionID(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case int64:
		return int(t)
	case string:
		n, _ := strconv.Atoi(t)
		return n
	default:
		return 0
	}
}
