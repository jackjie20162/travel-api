package logic

import (
	"context"
	"time"

	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/predicate"
	"gitee.com/meinongyihe/travel-rpc/ent/region"
	"gitee.com/meinongyihe/travel-rpc/ent/regiontranslation"
)

// baseLocale 区域基础名称语言：regions.name 存的就是这个语言的名称。
const baseLocale = "zh-CN"

// toRegionBase 将 ent.Region 映射为 API 类型（不含译文）。
func toRegionBase(r *ent.Region) types.Region {
	return types.Region{
		Id:        int64(r.ID),
		ParentId:  r.ParentID,
		Level:     r.Level,
		Code:      r.Code,
		Name:      r.Name,
		PoiId:     r.PoiID,
		PoiName:   r.PoiName,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
		Sort:      r.Sort,
		Status:    r.Status,
	}
}

// attachTranslations 批量加载区域译文并挂载到对应的 Region 上（管理端使用）。
func attachTranslations(ctx context.Context, db *ent.Client, items []types.Region) error {
	if len(items) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.Id)
	}
	rows, err := db.RegionTranslation.Query().
		Where(regiontranslation.RegionIDIn(ids...)).
		All(ctx)
	if err != nil {
		return err
	}
	byRegion := make(map[int64][]types.RegionTranslationItem, len(items))
	for _, t := range rows {
		byRegion[t.RegionID] = append(byRegion[t.RegionID], types.RegionTranslationItem{
			Locale: t.Locale,
			Name:   t.Name,
		})
	}
	for i := range items {
		items[i].Translations = byRegion[items[i].Id]
	}
	return nil
}

// resolveRegionNames 按请求语言解析区域名称：命中译文则覆盖基础名，否则回退 regions.name。
func resolveRegionNames(ctx context.Context, db *ent.Client, items []types.Region, locale string) error {
	if locale == "" || locale == baseLocale || len(items) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.Id)
	}
	rows, err := db.RegionTranslation.Query().
		Where(
			regiontranslation.RegionIDIn(ids...),
			regiontranslation.LocaleEQ(locale),
		).
		All(ctx)
	if err != nil {
		return err
	}
	nameByRegion := make(map[int64]string, len(rows))
	for _, t := range rows {
		if t.Name != "" {
			nameByRegion[t.RegionID] = t.Name
		}
	}
	for i := range items {
		if name, ok := nameByRegion[items[i].Id]; ok {
			items[i].Name = name
		}
	}
	return nil
}

// saveRegionTranslations 覆盖式写入区域译文：先删该区域全部译文，再按入参重建。
func saveRegionTranslations(ctx context.Context, db *ent.Client, regionID int64, items []types.RegionTranslationItem) error {
	if _, err := db.RegionTranslation.Delete().
		Where(regiontranslation.RegionIDEQ(regionID)).
		Exec(ctx); err != nil {
		return err
	}
	now := time.Now().Unix()
	for _, it := range items {
		if it.Locale == "" || it.Name == "" {
			continue
		}
		// 基础语言名称已存在 regions.name，无需重复落译文表
		if it.Locale == baseLocale {
			continue
		}
		if err := db.RegionTranslation.Create().
			SetRegionID(regionID).
			SetLocale(it.Locale).
			SetName(it.Name).
			SetSource("MANUAL").
			SetStatus("DONE").
			SetUpdatedAt(now).
			Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

// regionPredicates 依据列表请求构造区域查询条件。
// Level=COUNTRY 时限定 parent_id=0；Level=CITY 且带 ParentId 时按父级国家过滤。
func regionPredicates(req *types.RegionListReq) []predicate.Region {
	preds := make([]predicate.Region, 0, 4)
	switch req.Level {
	case "COUNTRY":
		preds = append(preds, region.ParentIDEQ(0))
	case "CITY":
		if req.ParentId > 0 {
			preds = append(preds, region.ParentIDEQ(req.ParentId))
		}
	default:
		if req.ParentId > 0 {
			preds = append(preds, region.ParentIDEQ(req.ParentId))
		}
	}
	if req.Level != "" {
		preds = append(preds, region.LevelEQ(req.Level))
	}
	if req.Keyword != "" {
		preds = append(preds, region.Or(region.NameContains(req.Keyword), region.CodeContains(req.Keyword)))
	}
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}
	preds = append(preds, region.StatusEQ(status))
	return preds
}
