package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListItineraryStopsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListItineraryStopsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListItineraryStopsLogic {
	return &ListItineraryStopsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListItineraryStopsLogic) ListItineraryStops() (*types.ItineraryStopListResp, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.ManagementClient.ListItineraryStops(ctx, &travel.ItineraryStopListRequest{ProductId: id})
	if err != nil {
		return nil, err
	}
	out := &types.ItineraryStopListResp{}
	for _, s := range resp.Items {
		out.Items = append(out.Items, toItineraryStop(s))
	}
	return out, nil
}
