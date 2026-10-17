package logic

import (
	"context"
	"encoding/json"

	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/attraction"
	"gitee.com/meinongyihe/travel-rpc/ent/itinerarystop"

	"github.com/zeromicro/go-zero/core/logx"
)

// fillMissingProductCovers 线路（产品）未设置主图时，回退取行程绑定景点的图片作为主图。
// 绑定关系记录在 itinerary_stops.type_params JSON 的 attractionId 字段中，
// 按行程顺序取第一个「有图」的景点：优先景点封面，其次景点图集第一张。
// 只补齐响应字段、不回写数据库，商户自行设置主图后自然优先使用商户的图。
func fillMissingProductCovers(ctx context.Context, db *ent.Client, items []types.Product) {
	missing := make([]int64, 0, len(items))
	for i := range items {
		if items[i].Id > 0 && items[i].CoverImage == "" {
			missing = append(missing, items[i].Id)
		}
	}
	if len(missing) == 0 {
		return
	}

	stops, err := db.ItineraryStop.Query().
		Where(itinerarystop.ProductIDIn(missing...)).
		Order(itinerarystop.ByProductID(), itinerarystop.BySequence(), itinerarystop.ByID()).
		All(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("load itinerary stops for cover fallback: %v", err)
		return
	}

	// productID → 按行程顺序去重后的景点 id
	plan := make(map[int64][]int, len(missing))
	seen := make(map[int64]map[int]bool, len(missing))
	all := make([]int, 0)
	for _, s := range stops {
		if s.TypeParams == "" {
			continue
		}
		var tp map[string]any
		if json.Unmarshal([]byte(s.TypeParams), &tp) != nil {
			continue
		}
		aid := parseAttractionID(tp["attractionId"])
		if aid <= 0 {
			continue
		}
		if seen[s.ProductID] == nil {
			seen[s.ProductID] = make(map[int]bool)
		}
		if seen[s.ProductID][aid] {
			continue
		}
		seen[s.ProductID][aid] = true
		plan[s.ProductID] = append(plan[s.ProductID], aid)
		all = append(all, aid)
	}
	if len(all) == 0 {
		return
	}

	rows, err := db.Attraction.Query().Where(attraction.IDIn(all...)).All(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("load attractions for cover fallback: %v", err)
		return
	}
	byID := make(map[int]*ent.Attraction, len(rows))
	for _, a := range rows {
		byID[a.ID] = a
	}

	for i := range items {
		if items[i].CoverImage != "" {
			continue
		}
		for _, aid := range plan[items[i].Id] {
			if img := attractionCoverCandidate(byID[aid]); img != "" {
				items[i].CoverImage = img
				break
			}
		}
	}
}

// attractionCoverCandidate 景点可用于线路主图的图片：封面优先，其次图集第一张。
func attractionCoverCandidate(a *ent.Attraction) string {
	if a == nil {
		return ""
	}
	if a.CoverImage != "" {
		return a.CoverImage
	}
	var urls []string
	if json.Unmarshal([]byte(a.Images), &urls) == nil {
		for _, u := range urls {
			if u != "" {
				return u
			}
		}
	}
	return ""
}
