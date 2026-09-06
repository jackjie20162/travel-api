package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListInventoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInventoryLogic {
	return &ListInventoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListInventoryLogic) ListInventory() (*types.InventoryListResp, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.ManagementClient.ListInventory(ctx, &travel.InventoryListRequest{PackageId: id})
	if err != nil {
		return nil, err
	}
	out := &types.InventoryListResp{}
	for _, i := range resp.Items {
		out.Items = append(out.Items, toInventoryItem(i))
	}
	return out, nil
}
