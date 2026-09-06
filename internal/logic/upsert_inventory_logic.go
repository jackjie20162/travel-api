package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpsertInventoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpsertInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertInventoryLogic {
	return &UpsertInventoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpsertInventoryLogic) UpsertInventory(req *types.InventoryUpsertReq) (*types.InventoryItem, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	item, err := l.svcCtx.ManagementClient.UpsertInventory(ctx, &travel.UpsertInventoryRequest{
		PackageId: id, Date: req.Date, TimeSlot: req.TimeSlot,
		Capacity: int32(req.Capacity), UnitPrice: req.UnitPrice,
		Currency: req.Currency, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	out := toInventoryItem(item)
	return &out, nil
}
