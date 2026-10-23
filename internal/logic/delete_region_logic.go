package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/attraction"
	"gitee.com/meinongyihe/travel-rpc/ent/region"
	"gitee.com/meinongyihe/travel-rpc/ent/regiontranslation"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRegionLogic {
	return &DeleteRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// DeleteRegion 删除区域：国家需先清空下级城市，城市需未被景点绑定，随后连带删除译文。
func (l *DeleteRegionLogic) DeleteRegion() (resp *types.OkResp, err error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	r, err := l.svcCtx.DB.Region.Query().Where(region.IDEQ(int(id))).Only(l.ctx)
	if err != nil {
		return nil, err
	}

	if r.Level == "COUNTRY" {
		childCount, err := l.svcCtx.DB.Region.Query().Where(region.ParentIDEQ(id)).Count(l.ctx)
		if err != nil {
			return nil, err
		}
		if childCount > 0 {
			return nil, errors.New("该国下仍有城市，无法删除")
		}
	} else {
		bound, err := l.svcCtx.DB.Attraction.Query().Where(attraction.CityIDEQ(id)).Exist(l.ctx)
		if err != nil {
			return nil, err
		}
		if bound {
			return nil, errors.New("该城市已被景点绑定，无法删除")
		}
	}

	if _, err := l.svcCtx.DB.RegionTranslation.Delete().
		Where(regiontranslation.RegionIDEQ(id)).Exec(l.ctx); err != nil {
		return nil, err
	}
	if err := l.svcCtx.DB.Region.DeleteOneID(int(id)).Exec(l.ctx); err != nil {
		return nil, err
	}
	return &types.OkResp{Ok: true}, nil
}
