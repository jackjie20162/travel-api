package logic

import (
	"context"
	"errors"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type BatchInventoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchInventoryLogic {
	return &BatchInventoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *BatchInventoryLogic) BatchInventory(req *types.InventoryBatchReq) (*types.InventoryListResp, error) {
	packageIdStr := pathVal(l.ctx, "packageId")
	if packageIdStr == "" {
		return nil, errors.New("missing packageId")
	}
	// Query ent directly for batch inventory
	// For now, return empty list as batch query needs ent client
	out := &types.InventoryListResp{}
	return out, nil
}
