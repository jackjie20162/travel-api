package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/region"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMerchantRegionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMerchantRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMerchantRegionsLogic {
	return &ListMerchantRegionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// ListMerchantRegions 管理端区域列表：返回基础名 + 全部译文，供后台编辑。
func (l *ListMerchantRegionsLogic) ListMerchantRegions(req *types.RegionListReq) (resp *types.RegionListResp, err error) {
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
	if err := attachTranslations(l.ctx, l.svcCtx.DB, items); err != nil {
		return nil, err
	}
	return &types.RegionListResp{Items: items, Total: int64(len(items))}, nil
}
