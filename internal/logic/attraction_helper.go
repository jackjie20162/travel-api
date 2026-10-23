package logic

import (
	"context"
	"encoding/json"
	"time"

	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/attractiontranslation"
)

// attachAttractionTranslations 批量加载景点译文并挂载到对应的 Attraction 上（管理端编辑用）。
// 译文表按 locale 存放 name/category/tags/address/description 的对应字段。
func attachAttractionTranslations(ctx context.Context, db *ent.Client, items []types.Attraction) error {
	if len(items) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.Id)
	}
	rows, err := db.AttractionTranslation.Query().
		Where(attractiontranslation.AttractionIDIn(ids...)).
		All(ctx)
	if err != nil {
		return err
	}
	byAttraction := make(map[int64][]types.AttractionTranslationItem, len(items))
	for _, t := range rows {
		byAttraction[t.AttractionID] = append(byAttraction[t.AttractionID], types.AttractionTranslationItem{
			Locale:      t.Locale,
			Name:        t.Name,
			Category:    t.Category,
			Tags:        t.Tags,
			Address:     t.Address,
			Description: t.Description,
		})
	}
	for i := range items {
		items[i].Translations = byAttraction[items[i].Id]
	}
	return nil
}

// saveAttractionTranslations 覆盖式写入景点译文：先删该景点全部译文，再按入参重建。
// 基础语言(zh-CN)内容存在 attractions 主表，译文表只保存 en-US/ar 等非基础语言。
func saveAttractionTranslations(ctx context.Context, db *ent.Client, attractionID int64, items []types.AttractionTranslationItem) error {
	if _, err := db.AttractionTranslation.Delete().
		Where(attractiontranslation.AttractionIDEQ(attractionID)).
		Exec(ctx); err != nil {
		return err
	}
	now := time.Now().Unix()
	for _, it := range items {
		if it.Locale == "" || it.Locale == baseLocale {
			continue
		}
		// 该 locale 无任何内容则跳过，避免写入空行
		if it.Name == "" && it.Category == "" && it.Tags == "" && it.Address == "" && it.Description == "" {
			continue
		}
		if err := db.AttractionTranslation.Create().
			SetAttractionID(attractionID).
			SetLocale(it.Locale).
			SetName(it.Name).
			SetCategory(it.Category).
			SetTags(it.Tags).
			SetAddress(it.Address).
			SetDescription(it.Description).
			SetSource("MANUAL").
			SetStatus("DONE").
			SetUpdatedAt(now).
			Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

// resolveAttractionLocale 按请求语言解析景点文本：命中译文则覆盖对应字段，否则回退主表基础语言(zh-CN)。
// 分类字段主表存的是枚举 code，locale 命中时用译文里的本地化显示名覆盖，供 C 端直接展示。
func resolveAttractionLocale(ctx context.Context, db *ent.Client, items []types.Attraction, locale string) error {
	if locale == "" || locale == baseLocale || len(items) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.Id)
	}
	rows, err := db.AttractionTranslation.Query().
		Where(
			attractiontranslation.AttractionIDIn(ids...),
			attractiontranslation.LocaleEQ(locale),
		).
		All(ctx)
	if err != nil {
		return err
	}
	byAttraction := make(map[int64]*ent.AttractionTranslation, len(rows))
	for _, t := range rows {
		byAttraction[t.AttractionID] = t
	}
	for i := range items {
		t, ok := byAttraction[items[i].Id]
		if !ok {
			continue
		}
		if t.Name != "" {
			items[i].Name = t.Name
		}
		if t.Category != "" {
			items[i].Category = t.Category
		}
		if t.Tags != "" {
			items[i].Tags = t.Tags
		}
		if t.Address != "" {
			items[i].Address = t.Address
		}
		if t.Description != "" {
			items[i].Description = t.Description
		}
		// 对外已按语言解析，译文明细无需返回
		items[i].Translations = nil
	}
	return nil
}

// resolveItineraryStopAttractionNames 按 locale 用景点译文覆盖行程节点的展示名称。
// 行程节点通过 type_params.attractionId 关联景点库；命中该 locale 的景点名称译文则
// 替换 title/poiName/locationName，未命中（无关联或无译文）则保留商户录入的原文。
// 基础语言(zh-CN)直接返回，保留中文原文。
func resolveItineraryStopAttractionNames(ctx context.Context, db *ent.Client, items []types.ItineraryStop, locale string) error {
	if locale == "" || locale == baseLocale || len(items) == 0 {
		return nil
	}
	stopAttraction := make(map[int]int64, len(items))
	ids := make([]int64, 0, len(items))
	for i := range items {
		if items[i].TypeParams == "" {
			continue
		}
		var tp map[string]any
		if json.Unmarshal([]byte(items[i].TypeParams), &tp) != nil {
			continue
		}
		aid := parseAttractionID(tp["attractionId"])
		if aid <= 0 {
			continue
		}
		stopAttraction[i] = int64(aid)
		ids = append(ids, int64(aid))
	}
	if len(ids) == 0 {
		return nil
	}
	rows, err := db.AttractionTranslation.Query().
		Where(
			attractiontranslation.AttractionIDIn(ids...),
			attractiontranslation.LocaleEQ(locale),
		).
		All(ctx)
	if err != nil {
		return err
	}
	nameByID := make(map[int64]string, len(rows))
	for _, t := range rows {
		if t.Name != "" {
			nameByID[t.AttractionID] = t.Name
		}
	}
	for i, aid := range stopAttraction {
		if name, ok := nameByID[aid]; ok && name != "" {
			items[i].Title = name
			items[i].PoiName = name
			items[i].LocationName = name
		}
	}
	return nil
}
