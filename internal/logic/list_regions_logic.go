package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/region"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRegionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionsLogic {
	return &ListRegionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// ListRegions 公开区域字典：返回国家或某国家下的城市，名称按 locale 解析。
func (l *ListRegionsLogic) ListRegions(req *types.RegionListReq) (resp *types.RegionListResp, err error) {
	rows, err := l.svcCtx.DB.Region.Query().
		Where(regionPredicates(req)...).
		Order(region.BySort(), region.ByID()).
		All(l.ctx)
	if err != nil {
		return nil, err
	}
	items := make([]types.Region, 0, len(rows))
	for _, r := range rows {
		items = append(items, toRegionBase(r))
	}
	if err := resolveRegionNames(l.ctx, l.svcCtx.DB, items, req.Locale); err != nil {
		return nil, err
	}
	return &types.RegionListResp{Items: items, Total: int64(len(items))}, nil
}
