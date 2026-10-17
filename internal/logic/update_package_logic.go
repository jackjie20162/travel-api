package logic

import (
	"context"
	"errors"
	"strings"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/productpackage"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePackageLogic {
	return &UpdatePackageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// UpdatePackage 按套餐 id 更新套餐字段。
// RPC 侧只有 CreatePackage（没有 UpdatePackage），此前这里误用 CreatePackage
// 并把套餐 id 当 product_id 传入，会在其它产品下生成脏套餐，因此改为 API 层直连 ent 更新
// （与景点等商户端接口一致的写法）。
func (l *UpdatePackageLogic) UpdatePackage(req *types.PackageReq) (*types.Package, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	tenantID, ok := middleware.TenantID(l.ctx)
	if !ok {
		return nil, errors.New("missing tenant context")
	}
	merchantID, ok := middleware.MerchantID(l.ctx)
	if !ok {
		return nil, errors.New("missing merchant context")
	}

	pkg, err := l.svcCtx.DB.ProductPackage.Query().
		Where(
			productpackage.IDEQ(int(id)),
			productpackage.TenantIDEQ(tenantID),
			productpackage.MerchantIDEQ(merchantID),
		).
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("package not found")
		}
		return nil, err
	}

	upd := pkg.Update()
	if code := strings.TrimSpace(req.Code); code != "" && code != pkg.Code {
		upd.SetCode(code)
	}
	if name := strings.TrimSpace(req.Name); name != "" {
		upd.SetName(name)
	}
	if req.PricingMode != "" {
		upd.SetPricingMode(req.PricingMode)
	}
	if req.InventoryMode != "" {
		upd.SetInventoryMode(req.InventoryMode)
	}
	if req.MinOrderQty > 0 {
		upd.SetMinOrderQty(req.MinOrderQty)
	}
	if req.SellCurrency != "" {
		upd.SetSellCurrency(strings.ToUpper(req.SellCurrency))
	}
	if req.CostCurrency != "" {
		upd.SetCostCurrency(strings.ToUpper(req.CostCurrency))
	}
	if req.GroupPrices != "" {
		upd.SetGroupPrices(req.GroupPrices)
	}
	if req.TierPrices != "" {
		upd.SetTierPrices(req.TierPrices)
	}

	saved, err := upd.Save(l.ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errors.New("套餐编码在当前产品下已存在")
		}
		return nil, err
	}
	return entToPackage(saved), nil
}

// entToPackage 将 ent.ProductPackage 映射为 API 响应类型。
func entToPackage(p *ent.ProductPackage) *types.Package {
	return &types.Package{
		Id: int64(p.ID), ProductId: p.ProductID, Code: p.Code, Name: p.Name,
		Status: p.Status, PricingMode: p.PricingMode, InventoryMode: p.InventoryMode,
		MinOrderQty: p.MinOrderQty, SellCurrency: p.SellCurrency, CostCurrency: p.CostCurrency,
		GroupPrices: p.GroupPrices, TierPrices: p.TierPrices,
	}
}
