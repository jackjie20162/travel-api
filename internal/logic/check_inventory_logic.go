package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckInventoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInventoryLogic {
	return &CheckInventoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CheckInventoryLogic) CheckInventory(req *types.InventoryCheckReq) (*types.InventoryCheckResp, error) {
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.InventoryClient.Check(ctx, &travel.InventoryRequest{
		PackageId: req.PackageId, Date: req.Date,
		TimeSlot: req.TimeSlot, Quantity: int32(req.Quantity),
	})
	if err != nil {
		return nil, err
	}
	return &types.InventoryCheckResp{
		Available: resp.Available, Remaining: resp.Remaining,
		UnitPrice: resp.UnitPrice, Currency: resp.Currency,
	}, nil
}
