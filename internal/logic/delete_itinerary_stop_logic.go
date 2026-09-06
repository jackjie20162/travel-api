package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteItineraryStopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteItineraryStopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteItineraryStopLogic {
	return &DeleteItineraryStopLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteItineraryStopLogic) DeleteItineraryStop() (*types.OkResp, error) {
	stopId, err := pathInt64(l.ctx, "stopId")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ManagementClient.DeleteItineraryStop(ctx, &travel.DeleteItineraryStopRequest{Id: stopId})
	if err != nil {
		return nil, err
	}
	return &types.OkResp{Ok: true}, nil
}
