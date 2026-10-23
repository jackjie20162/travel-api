package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/i18n"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductItineraryStopsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductItineraryStopsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductItineraryStopsLogic {
	return &ListProductItineraryStopsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *ListProductItineraryStopsLogic) ListProductItineraryStops() (resp *types.ItineraryStopListResp, err error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.CatalogClient.ListItineraryStops(ctx, &travel.ItineraryStopListRequest{ProductId: id})
	if err != nil {
		return nil, err
	}
	out := &types.ItineraryStopListResp{}
	for _, s := range result.Items {
		out.Items = append(out.Items, toItineraryStop(s))
	}
	// 行程节点面向 C 端展示：关联景点的名称按请求语言解析（中/英/阿）。
	if err := resolveItineraryStopAttractionNames(l.ctx, l.svcCtx.DB, out.Items, i18n.LocaleFrom(l.ctx)); err != nil {
		return nil, err
	}
	return out, nil
}
