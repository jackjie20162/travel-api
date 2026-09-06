package logic

import (
	"context"
	"errors"
	"strconv"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/inventory"
	"gitee.com/meinongyihe/travel-rpc/ent/predicate"

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
	packageID, err := strconv.ParseInt(packageIdStr, 10, 64)
	if err != nil || packageID <= 0 {
		return nil, errors.New("invalid packageId")
	}

	// Build ent query: filter by packageId + date range, optionally by merchantID
	preds := []predicate.Inventory{
		inventory.PackageIDEQ(packageID),
		inventory.ServiceDateGTE(req.StartDate),
		inventory.ServiceDateLTE(req.EndDate),
	}
	if merchantID, ok := middleware.MerchantID(l.ctx); ok && merchantID > 0 {
		preds = append(preds, inventory.MerchantIDEQ(merchantID))
	}

	rows, err := l.svcCtx.DB.Inventory.Query().Where(preds...).Order(ent.Asc(inventory.FieldServiceDate)).All(l.ctx)
	if err != nil {
		return nil, err
	}

	out := &types.InventoryListResp{}
	for _, r := range rows {
		out.Items = append(out.Items, types.InventoryItem{
			Id:            int64(r.ID),
			PackageId:     r.PackageID,
			Date:          r.ServiceDate,
			TimeSlot:      r.TimeSlot,
			Capacity:      int32(r.Capacity),
			Reserved:      int32(r.Reserved),
			UnitPrice:     r.UnitPrice,
			Currency:      r.Currency,
			Status:        r.Status,
			InventoryMode: r.InventoryMode,
			TotalCapacity: int32(r.TotalCapacity),
			IsOpen:        r.IsOpen,
		})
	}
	return out, nil
}
