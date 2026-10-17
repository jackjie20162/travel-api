package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePackageLogic {
	return &CreatePackageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreatePackageLogic) CreatePackage(req *types.PackageReq) (*types.Package, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	p, err := l.svcCtx.ManagementClient.CreatePackage(ctx, &travel.CreatePackageRequest{
		ProductId: id, Code: req.Code, Name: req.Name,
		// 报价/库存/币种等扩展字段一并下传，RPC 侧有默认值兼容空值
		PricingMode: req.PricingMode, InventoryMode: req.InventoryMode,
		MinOrderQty: int32(req.MinOrderQty), SellCurrency: req.SellCurrency,
		CostCurrency: req.CostCurrency, GroupPrices: req.GroupPrices,
		TierPrices: req.TierPrices,
	})
	if err != nil {
		return nil, err
	}
	out := toPackage(p)
	return &out, nil
}
