package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/region"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRegionLogic {
	return &UpdateRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// UpdateRegion 更新国家/城市基础信息并覆盖式重建译文。
func (l *UpdateRegionLogic) UpdateRegion(req *types.RegionReq) (resp *types.Region, err error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	exist, err := l.svcCtx.DB.Region.Query().Where(region.IDEQ(int(id))).Exist(l.ctx)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("region not found")
	}

	parentID := req.ParentId
	if req.Level == "COUNTRY" {
		parentID = 0
	}
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	r, err := l.svcCtx.DB.Region.UpdateOneID(int(id)).
		SetParentID(parentID).
		SetLevel(req.Level).
		SetCode(req.Code).
		SetName(req.Name).
		SetPoiID(req.PoiId).
		SetPoiName(req.PoiName).
		SetLatitude(req.Latitude).
		SetLongitude(req.Longitude).
		SetSort(req.Sort).
		SetStatus(status).
		Save(l.ctx)
	if err != nil {
		return nil, err
	}

	if err := saveRegionTranslations(l.ctx, l.svcCtx.DB, id, req.Translations); err != nil {
		return nil, err
	}

	out := toRegionBase(r)
	items := []types.Region{out}
	if err := attachTranslations(l.ctx, l.svcCtx.DB, items); err != nil {
		return nil, err
	}
	out = items[0]
	return &out, nil
}
