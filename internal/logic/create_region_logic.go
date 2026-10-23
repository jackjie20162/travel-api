package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRegionLogic {
	return &CreateRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// CreateRegion 新建国家或城市；国家 parent_id 固定为 0。
func (l *CreateRegionLogic) CreateRegion(req *types.RegionReq) (resp *types.Region, err error) {
	parentID := req.ParentId
	if req.Level == "COUNTRY" {
		parentID = 0
	}
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	r, err := l.svcCtx.DB.Region.Create().
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

	if err := saveRegionTranslations(l.ctx, l.svcCtx.DB, int64(r.ID), req.Translations); err != nil {
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
